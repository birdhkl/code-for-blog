package arch_test

import (
	"context"
	"fmt"
	"software-arch/arch"
	"sync"
	"testing"
)

const msg = "Hello World"

var eventOnInput = context.WithValue(context.Background(), arch.GLabelEventType, arch.EventOnInput)
var eventOnScrollAdded = context.WithValue(context.Background(), arch.GLabelEventType, arch.EventOnScrollAdded)
var eventOnTextChanged = context.WithValue(context.Background(), arch.GLabelEventType, arch.EventOnTextChanged)

type ClientViewDev struct {
	rw     sync.RWMutex
	buffer []byte
}

func (dev *ClientViewDev) Write(p []byte) (n int, err error) {
	dev.rw.Lock()
	defer dev.rw.Unlock()
	dev.buffer = append(dev.buffer, p...)
	return len(p), nil
}

func (dev *ClientViewDev) Read(p []byte) (n int, err error) {
	dev.rw.Lock()
	defer dev.rw.Unlock()
	if dev.buffer == nil {
		return 0, fmt.Errorf("buffer nil")
	}
	nLen := len(p)
	if nLen > len(dev.buffer) {
		nLen = len(dev.buffer)
	}
	for index := 0; index < nLen; index++ {
		p[index] = dev.buffer[index]
	}
	dev.buffer = dev.buffer[nLen:]
	return nLen, nil
}

func (w *ClientViewDev) ToString() string {
	w.rw.RLock()
	defer w.rw.RUnlock()
	return string(w.buffer)
}

func (w *ClientViewDev) ReadString() (string, error) {
	buffer := make([]byte, 1024)
	n, err := w.Read(buffer)
	if err != nil {
		return "", err
	}
	res := string(buffer[:n])
	return res, nil
}

func TestLayer(t *testing.T) {
	dev := &ClientViewDev{}
	low := arch.NewLowLayer(dev)
	mid := arch.NewMiddleLayer(low)
	top := arch.NewTopLayer(mid)
	err := top.Do(context.WithValue(context.Background(), arch.GLabelContent, msg))
	if err != nil {
		t.Error(err)
	}
	wantStr := fmt.Sprintf("#LowLayer[#MiddleLayer[#TopLayer[%s]]]", msg)
	if dev.ToString() != wantStr {
		t.Errorf("Get(%s) != %s", dev.ToString(), msg)
	}
}
