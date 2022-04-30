package arch_test

import (
	"context"
	"software-arch/arch/mka/kernel"
	"sync"
	"testing"
)

func TestMicrokernel(t *testing.T) {
	// 先编译mka, cd mka && make build-bundle
	const bundlePath = "./mka/plugins"
	const bundleName = "hello_world"
	const serviceName = "Service"

	bundleContext := kernel.NewDefaultBundleContext(bundlePath, context.TODO())
	if err := bundleContext.InstallBundle(bundleName); err != nil {
		t.Error(err)
		return
	}
	installBundle, err := bundleContext.GetBundle(bundleName)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err := bundleContext.Stop(); err != nil {
			t.Error(err)
		}
		if _, err := bundleContext.GetBundle(bundleName); err == nil {
			t.Errorf("bundle %s not unregister", bundleName)
		}
		if len(installBundle.GetBundleServices()) != 0 {
			t.Errorf("service not unregister, %v", installBundle.GetBundleServices())
		}
	}()

	testCases := []struct {
		funcName string
		msg      string
		result   string
	}{
		{
			funcName: "test1",
			msg:      "msg1",
			result:   "test1#Hello World#msg1",
		},
		{
			funcName: "test2",
			msg:      "msg2",
			result:   "test2#Hello World#msg2",
		},
		{
			funcName: "test3",
			msg:      "msg3",
			result:   "test3#Hello World#msg3",
		},
		{
			funcName: "test4",
			msg:      "msg4",
			result:   "test4#Hello World#msg4",
		},
	}
	g := sync.WaitGroup{}
	for _, testCase := range testCases {
		g.Add(1)
		go func(funcName string, sendmsg string, result string) {
			defer g.Done()
			srv, ok := installBundle.GetBundleServices()[serviceName]
			if !ok {
				t.Errorf("srv %s not register", serviceName)
				return
			}
			srvRef, err := bundleContext.GetServiceReference(serviceName)
			if err != nil {
				t.Error(err)
				return
			}
			msg := kernel.NewDefaultMessage(funcName, sendmsg)
			res := <-srvRef.Send(msg)
			if res.Value() != result {
				t.Errorf("serviceRef, want %s, but %v", result, res.Value())
			}
			res = srv.Recv(msg)
			if res.Value() != result {
				t.Errorf("service, want %s, but %v", result, res.Value())
			}
		}(testCase.funcName, testCase.msg, testCase.result)
	}
	g.Wait()
}
