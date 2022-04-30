package main

import (
	"fmt"
	"software-arch/arch/mka/bundle"
	"reflect"
)

type Message struct {
	value    string
	funcName string
}

func (m *Message) Error() error {
	return nil
}

func (m *Message) Value() interface{} {
	return m.funcName + "#" + "Hello World" + "#" + m.value
}

type HelloWorldService struct {
}

func (s *HelloWorldService) Recv(msg bundle.BundleServiceMessage) bundle.Message {
	value, ok := msg.GetMessage().(string)
	if !ok {
		fmt.Println("GetMessage ", msg.GetMessage(), "type", reflect.TypeOf(msg.GetMessage()))
	}
	return &Message{value: value, funcName: msg.GetFunctionName()}
}

type ServiceActivator struct {
}

func (activator *ServiceActivator) Start(ctx bundle.BundleContext) {
	if err := ctx.RegisterService("Service", &HelloWorldService{}); err != nil {
		fmt.Println("Start ", err)
	}

}
func (activator *ServiceActivator) Stop(ctx bundle.BundleContext) {
	if err := ctx.UnregisterService("Service"); err != nil {
		fmt.Println("Stop", err)
	}
}

var Activator ServiceActivator
