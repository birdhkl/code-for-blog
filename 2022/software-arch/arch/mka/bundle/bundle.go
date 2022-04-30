package bundle

import (
	"fmt"
	"plugin"
	"sync"
)

// Bundle is the abstract layer for golang plugin
type Bundle struct {
	URL             string `json:"url"`
	Version         string `json:"version"`
	Name            string `json:"name"`
	Desc            string `json:"description"`
	SymbolActivator string `json:"activator"`
	services        map[string]BundleService
	mux             sync.Mutex
}

func NewBundle() *Bundle {
	return &Bundle{
		URL:             "",
		Version:         "",
		Name:            "",
		Desc:            "",
		SymbolActivator: "",
		services:        map[string]BundleService{},
		mux:             sync.Mutex{},
	}
}

func (b *Bundle) GetBundleUrl() string {
	return b.URL
}

func (b *Bundle) GetBundleVersion() string {
	return b.Version
}

func (b *Bundle) GetBundleName() string {
	return b.Name
}

func (b *Bundle) GetBundleDescription() string {
	return b.Desc
}

func (b *Bundle) RegisterService(name string, srv BundleService) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.services[name] = srv
}
func (b *Bundle) UnregisterService(name string) {
	b.mux.Lock()
	defer b.mux.Unlock()
	delete(b.services, name)
}

func (b *Bundle) GetBundleServices() map[string]BundleService {
	b.mux.Lock()
	defer b.mux.Unlock()
	services := make(map[string]BundleService, len(b.services))
	for name, srv := range b.services {
		services[name] = srv
	}
	return services
}

func (b *Bundle) GetBundleActivator() (BundleActivator, error) {
	symbol, err := b.lookUp(b.SymbolActivator)
	if err != nil {
		return nil, fmt.Errorf("Symbol %s Not Exist, err=%v", b.SymbolActivator, err)
	}
	activator, ok := symbol.(BundleActivator)
	if !ok {
		return nil, fmt.Errorf("No %s in Bundle[%s]", b.SymbolActivator, b.Name)
	}
	return activator, nil
}

func (b *Bundle) lookUp(symbol string) (plugin.Symbol, error) {
	url := b.GetBundleUrl()
	plug, err := plugin.Open(url)
	if err != nil {
		fmt.Println("plugin.Open err=", err)
		return nil, err
	}
	return plug.Lookup(b.SymbolActivator)
}
