package arch

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ServiceGenerator func(k ServiceProtocolKind) (EnterpriseService, error)
type ServiceContent map[string]interface{}

type ServiceProtocolKind uint32

const (
	_ ServiceProtocolKind = iota
	ServiceJSON
	ServiceRPC
)

type EnterpriseServiceBus interface {
	PostEventAsBroken(srcName, dstName, event string, content ServiceContent) (ServiceContent, error)
	PostEventAsMQ(srcName, onEvent string, dstName, event string, content ServiceContent) error
	Ping(srcName string, k ServiceProtocolKind) error
}

type EnterpriseService interface {
	Pong() error
	Name() string
	Kind() ServiceProtocolKind
	RecvEvent(name string, content ServiceContent) (ServiceContent, error)
}

type EnterpriseServiceBusService interface {
	EnterpriseServiceBus
	Start(ctx context.Context)
	Stop() <-chan error
}

func NewMemEnterpriseServiceBus(g ServiceGenerator) EnterpriseServiceBusService {
	return &MemBus{
		wg:                sync.WaitGroup{},
		srvName2Srv:       nil,
		srvName2Dead:      nil,
		rwmu:              sync.RWMutex{},
		mq:                nil,
		status:            0,
		stop:              nil,
		ctx:               nil,
		kind2SrvGenerator: g,
	}
}

type MemBus struct {
	wg                sync.WaitGroup
	srvName2Srv       map[string]EnterpriseService
	srvName2Dead      map[string]int64
	rwmu              sync.RWMutex
	mq                chan func() error
	status            uint32
	stop              func() <-chan error
	ctx               context.Context
	kind2SrvGenerator ServiceGenerator
}

func (m *MemBus) Stop() <-chan error {
	if atomic.LoadUint32(&m.status) != 1 {
		ch := make(chan error)
		defer close(ch)
		ch <- fmt.Errorf("Bus Not Useful")
		return ch
	}
	return m.stop()
}

func (m *MemBus) Start(ctx context.Context) {
	if !atomic.CompareAndSwapUint32(&m.status, 0, 1) {
		return
	}
	m.rwmu.Lock()
	defer m.rwmu.Unlock()
	ctx, cancel := context.WithCancel(ctx)
	m.ctx = ctx
	m.srvName2Srv = make(map[string]EnterpriseService)
	m.srvName2Dead = make(map[string]int64)
	m.mq = make(chan func() error, 100)

	m.stop = func() <-chan error {
		ch := make(chan error)
		go func() {
			defer close(ch)
			<-time.After(time.Millisecond * 10)
			if !atomic.CompareAndSwapUint32(&m.status, 1, 2) {
				ch <- fmt.Errorf("ebs not run")
				return
			}
			cancel()
			m.wg.Wait()
			m.rwmu.Lock()
			defer m.rwmu.Unlock()
			m.srvName2Srv = nil
			m.srvName2Dead = nil
			m.mq = nil
			m.ctx = nil
		}()
		return ch
	}

	go func() {
		defer func() {
			m.wg.Done()
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		m.wg.Add(1)
		for {
			const checkIntervalSec = 60
			const serviceDeadSec = 20
			select {
			case <-m.ctx.Done():
				return
			case cb := <-m.mq:
				go func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Println("Recovered in f", r)
						}
					}()
					if err := cb(); err != nil {
						fmt.Println("message err", err)
					}
				}()

			case <-time.After(time.Second * checkIntervalSec):
				go func() {
					checkEpochSec := time.Now().Unix()
					m.rwmu.RLock()
					deadServices := make([]string, 0, len(m.srvName2Srv))
					for service, pingEpochSec := range m.srvName2Dead {
						if checkEpochSec-pingEpochSec > serviceDeadSec {
							deadServices = append(deadServices, service)
						}
					}
					m.rwmu.RUnlock()

					m.rwmu.Lock()
					defer m.rwmu.Unlock()
					for _, service := range deadServices {
						delete(m.srvName2Dead, service)
						delete(m.srvName2Srv, service)
					}
				}()
			}
		}
	}()
}

func (m *MemBus) Ping(serviceName string, k ServiceProtocolKind) error {
	m.rwmu.RLock()
	p, ok := m.srvName2Srv[serviceName]
	m.rwmu.RUnlock()
	if !ok || p.Kind() != k {
		gp, err := m.kind2SrvGenerator(k)
		if err != nil {
			return err
		}
		p = gp
	}
	m.rwmu.Lock()
	if !ok {
		m.srvName2Srv[serviceName] = p
	}
	m.srvName2Dead[serviceName] = time.Now().Unix()
	m.rwmu.Unlock()

	if err := p.Pong(); err != nil {
		return err
	}

	return nil
}

func (m *MemBus) PostEventAsBroken(srcName, dstName, event string, content ServiceContent) (ServiceContent, error) {
	if atomic.LoadUint32(&m.status) != 1 {
		return nil, fmt.Errorf("bus is to be closed")
	}
	m.rwmu.RLock()
	p, ok := m.srvName2Srv[dstName]
	m.rwmu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%s not exist", srcName)
	}
	m.wg.Add(1)
	defer m.wg.Done()
	return p.RecvEvent(event, content)
}

func (m *MemBus) PostEventAsMQ(srcName, onEvent string, dstName, event string, content ServiceContent) error {
	if atomic.LoadUint32(&m.status) != 1 {
		return fmt.Errorf("bus is to be closed")
	}
	m.wg.Add(1)
	m.mq <- func() error {
		defer m.wg.Done()
		content, err := m.PostEventAsBroken(srcName, dstName, event, content)
		if err != nil {
			content = make(ServiceContent)
			content["Error"] = err.Error()
		}
		_, err = m.PostEventAsBroken(dstName, srcName, onEvent, content)
		return err
	}
	return nil
}
