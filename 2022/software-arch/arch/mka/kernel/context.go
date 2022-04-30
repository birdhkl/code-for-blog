package kernel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"software-arch/arch/mka/bundle"
	"sync"
	"time"
)

const (
	_ = iota
	ContextStart
	ContextStop
)

func NewBundleServiceProxy(srv bundle.BundleService, cancel func()) *BundleServiceProxy {
	return &BundleServiceProxy{
		BundleService: srv,
		cancel:        cancel,
		messageQueue:  make(chan bundle.BundleServiceMessage, 10),
	}
}

type BundleServiceProxy struct {
	bundle.BundleService
	cancel       func()
	messageQueue chan bundle.BundleServiceMessage
}

func (p *BundleServiceProxy) Stop() {
	p.cancel()
	close(p.messageQueue)
}

func (p *BundleServiceProxy) Queue() chan bundle.BundleServiceMessage {
	return p.messageQueue
}

type DefaultBundleContext struct {
	bundlePath   string
	rwMux        sync.RWMutex
	bundles      map[string]*bundle.Bundle
	services     map[string]*BundleServiceProxy
	context      context.Context
	state        int64
	serviceGroup sync.WaitGroup
}

func NewDefaultBundleContext(path string, ctx context.Context) bundle.BundleContext {
	return &DefaultBundleContext{
		bundlePath:   path,
		rwMux:        sync.RWMutex{},
		bundles:      map[string]*bundle.Bundle{},
		services:     map[string]*BundleServiceProxy{},
		context:      ctx,
		serviceGroup: sync.WaitGroup{},
		state:        ContextStart,
	}
}

func (ctx *DefaultBundleContext) Stop() error {
	// change state
	ctx.rwMux.Lock()
	if ctx.state != ContextStart {
		ctx.rwMux.Unlock()
		return fmt.Errorf("context not running")
	}
	ctx.rwMux.Unlock()
	// uninstall all bundles
	bundles, err := ctx.GetBundles()
	if err != nil {
		return err
	}
	for _, b := range bundles {
		if err := ctx.UninstallBundle(b.GetBundleName()); err != nil {
			fmt.Printf("err = %v\n", err)
		}
	}
	c := make(chan bool)
	go func() {
		// 等待所有service退出
		ctx.serviceGroup.Wait()
		close(c)
	}()
	select {
	case <-c:
	case <-time.After(time.Second * 5):
		fmt.Println("service not all quit")
	}
	return nil
}

func (ctx *DefaultBundleContext) GetBundles() ([]*bundle.Bundle, error) {
	ctx.rwMux.RLock()
	defer ctx.rwMux.RUnlock()
	bundles := make([]*bundle.Bundle, 0, len(ctx.bundles))
	for _, b := range ctx.bundles {
		bundles = append(bundles, b)
	}
	return bundles, nil
}

func (ctx *DefaultBundleContext) GetBundle(bundleName string) (*bundle.Bundle, error) {
	ctx.rwMux.RLock()
	defer ctx.rwMux.RUnlock()
	b, ok := ctx.bundles[bundleName]
	if !ok {
		return nil, fmt.Errorf("bundle %s not exist", bundleName)
	}
	return b, nil
}

func (ctx *DefaultBundleContext) InstallBundle(bundleName string) error {
	bundleConfigPath := path.Join(ctx.bundlePath, bundleName+".json")
	f, err := os.Open(bundleConfigPath)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	b := bundle.NewBundle()
	if err := json.Unmarshal(content, &b); err != nil {
		return err
	}
	if b.GetBundleName() != bundleName {
		return fmt.Errorf("want bundle %s, but real name %s", bundleName, b.GetBundleName())
	}
	activator, err := b.GetBundleActivator()
	if err != nil {
		return err
	}
	ctx.rwMux.Lock()
	if ctx.state != ContextStart {
		ctx.rwMux.Unlock()
		return fmt.Errorf("context not running")
	}
	if _, ok := ctx.bundles[b.GetBundleName()]; ok {
		ctx.rwMux.Unlock()
		return fmt.Errorf("%s repeat installed", b.GetBundleName())
	}
	ctx.bundles[b.GetBundleName()] = b
	ctx.rwMux.Unlock()
	activator.Start(NewSpecificBundleContextMiddleware(ctx, b))
	return nil
}

func (ctx *DefaultBundleContext) UninstallBundle(bundleName string) error {
	ctx.rwMux.Lock()
	b, ok := ctx.bundles[bundleName]
	if !ok {
		ctx.rwMux.Unlock()
		return fmt.Errorf("bundle %s not install", bundleName)
	}
	delete(ctx.bundles, bundleName)
	ctx.rwMux.Unlock()
	activator, err := b.GetBundleActivator()
	if err != nil {
		return err
	}
	activator.Stop(NewSpecificBundleContextMiddleware(ctx, b))
	return nil
}

func (ctx *DefaultBundleContext) RegisterService(serviceName string, srv bundle.BundleService) error {
	c, cancel := context.WithCancel(ctx.context)
	p := NewBundleServiceProxy(srv, cancel)
	ctx.rwMux.Lock()
	if ctx.state != ContextStart {
		ctx.rwMux.Unlock()
		return fmt.Errorf("context not running")
	}
	if _, ok := ctx.services[serviceName]; ok {
		ctx.rwMux.Unlock()
		return fmt.Errorf("service %s repeat", serviceName)
	}
	ctx.services[serviceName] = p
	ctx.rwMux.Unlock()
	ctx.serviceGroup.Add(1)
	go func() {
		defer func() {
			defer ctx.serviceGroup.Done()
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		fmt.Printf("service %s start\n", serviceName)
		for {
			select {
			case e, ok := <-p.Queue():
				if !ok {
					fmt.Printf("service %s queue closed\n", serviceName)
					break
				}
				msg := p.Recv(e)
				if resE, ok := e.(*messageWithResult); ok {
					resE.GetResultChan() <- msg
					close(resE.GetResultChan())
				}
				fmt.Println("recv:", e, "msg:", msg, "type:", reflect.TypeOf(e))
			case <-c.Done():
				fmt.Printf("service %s done\n", serviceName)
				return
			}
		}
	}()
	return nil
}

func (ctx *DefaultBundleContext) UnregisterService(serviceName string) error {
	ctx.rwMux.Lock()
	p, ok := ctx.services[serviceName]
	if !ok {
		ctx.rwMux.Unlock()
		return fmt.Errorf("%s not register", serviceName)
	}
	delete(ctx.services, serviceName)
	ctx.rwMux.Unlock()
	p.Stop()
	return nil
}

func (ctx *DefaultBundleContext) GetServiceReference(srvName string) (bundle.BundleServiceReference, error) {
	ctx.rwMux.RLock()
	defer ctx.rwMux.RUnlock()
	p, ok := ctx.services[srvName]
	if !ok {
		return nil, fmt.Errorf("%s not exist", srvName)
	}
	return NewDefaultBundleSrvRef(p.Queue()), nil
}

type SpecificBundleContext struct {
	bundle.BundleContext
	bundle *bundle.Bundle
}

func NewSpecificBundleContextMiddleware(ctx bundle.BundleContext, b *bundle.Bundle) bundle.BundleContext {
	return &SpecificBundleContext{
		BundleContext: ctx,
		bundle:        b,
	}
}

func (sctx *SpecificBundleContext) RegisterService(name string, srv bundle.BundleService) error {
	if err := sctx.BundleContext.RegisterService(name, srv); err != nil {
		return err
	}
	sctx.bundle.RegisterService(name, srv)
	return nil
}

func (sctx *SpecificBundleContext) UnregisterService(name string) error {
	if err := sctx.BundleContext.UnregisterService(name); err != nil {
		return err
	}
	sctx.bundle.UnregisterService(name)
	return nil
}
