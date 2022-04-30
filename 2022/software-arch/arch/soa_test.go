package arch_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"software-arch/arch"
	"sync"
	"testing"
)

type Service interface {
	Recv([]byte) ([]byte, error)
}

func NewServiceImpl() Service {
	return &ServiceImpl{}
}

type ServiceImpl struct {
	mu     sync.Mutex
	events []string
}

func (srv *ServiceImpl) Recv(event []byte) ([]byte, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.events = append(srv.events, string(event))
	return event, nil
}

type RpcEnterpriseService struct {
	srv  Service
	name string
}

func (srv *RpcEnterpriseService) Pong() error {
	return nil
}

func (srv *RpcEnterpriseService) Name() string {
	return srv.name
}

func (srv *RpcEnterpriseService) Kind() arch.ServiceProtocolKind {
	return arch.ServiceRPC
}

func (srv *RpcEnterpriseService) RecvEvent(event string,
	content arch.ServiceContent) (arch.ServiceContent, error) {
	bytes, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	eventBytes := []byte(event)
	if len(eventBytes) > 20 {
		return nil, fmt.Errorf("Rpc Name %s Too Long", event)
	}
	for len(eventBytes) < 20 {
		eventBytes = append(eventBytes, ' ')
	}
	respBytes, err := srv.srv.Recv(append(eventBytes, bytes...))
	if err != nil {
		return nil, err
	}
	resp := make(arch.ServiceContent)
	if err := json.Unmarshal(respBytes[20:], &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type JsonEnterpriseService struct {
	srv  Service
	name string
}

func (srv *JsonEnterpriseService) Pong() error {
	return nil
}

func (srv *JsonEnterpriseService) Name() string {
	return srv.name
}

func (srv *JsonEnterpriseService) Kind() arch.ServiceProtocolKind {
	return arch.ServiceJSON
}

func (srv *JsonEnterpriseService) RecvEvent(event string,
	content arch.ServiceContent) (arch.ServiceContent, error) {
	serviceBody := make(arch.ServiceContent, len(content))
	for k, v := range content {
		serviceBody[k] = v
	}
	serviceBody["_method"] = event

	bytes, err := json.Marshal(serviceBody)
	if err != nil {
		return nil, err
	}
	respBytes, err := srv.srv.Recv(bytes)
	if err != nil {
		return nil, err
	}
	resp := make(arch.ServiceContent)
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	delete(resp, "_method")
	return resp, nil
}

func TestSOA(t *testing.T) {
	jsonSrv := &JsonEnterpriseService{name: "json", srv: NewServiceImpl()}
	rpcSrv := &RpcEnterpriseService{name: "rpc", srv: NewServiceImpl()}
	g := func(k arch.ServiceProtocolKind) (arch.EnterpriseService, error) {
		switch k {
		case arch.ServiceJSON:
			return jsonSrv, nil
		case arch.ServiceRPC:
			return rpcSrv, nil
		default:
			return nil, fmt.Errorf("No protocol %v", k)
		}
	}
	// prepare service bus
	bus := arch.NewMemEnterpriseServiceBus(g)
	bus.Start(context.Background())
	if err := bus.Ping("json", arch.ServiceJSON); err != nil {
		t.Error(err)
		return
	}
	if err := bus.Ping("rpc", arch.ServiceRPC); err != nil {
		t.Error(err)
		return
	}
	// post json->rpc
	c1 := arch.ServiceContent{"123": "456"}
	resp, err := bus.PostEventAsBroken("json", "rpc", "hello", c1)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(&resp, &c1) {
		t.Errorf("%v != %v", resp, c1)
	}
	// post rpc->rpc
	c2 := arch.ServiceContent{"456": "789"}
	resp, err = bus.PostEventAsBroken("rpc", "json", "world", c2)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(&resp, &c2) {
		t.Errorf("%v != %v", resp, c2)
	}
	// check rpcSrv
	events := rpcSrv.srv.(*ServiceImpl).events
	if len(events) != 1 || events[0] != "hello               {\"123\":\"456\"}" {
		t.Errorf("events %v, not size one", events)
	}
	// check jsonSrv
	events = jsonSrv.srv.(*ServiceImpl).events
	if len(events) != 1 || events[0] != "{\"456\":\"789\",\"_method\":\"world\"}" {
		t.Errorf("events %v, not size one", events)
	}
	c3 := arch.ServiceContent{"789": "91011"}
	if err := bus.PostEventAsMQ("json", "onTony", "rpc", "tony", c3); err != nil {
		t.Error(err)
		return
	}
	<-bus.Stop()
	// check rpcSrv
	events = rpcSrv.srv.(*ServiceImpl).events
	if len(events) != 2 || events[1] != "tony                {\"789\":\"91011\"}" {
		t.Errorf("events %v, not size one", events)
	}
	// check jsonSrv
	events = jsonSrv.srv.(*ServiceImpl).events
	if len(events) != 2 || events[1] != "{\"789\":\"91011\",\"_method\":\"onTony\"}" {
		t.Errorf("events %v, not size one", events)
	}

}
