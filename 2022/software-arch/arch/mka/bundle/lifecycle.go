package bundle

const (
	_ int64 = iota
	BundleStart
	BundleStop
)

// BundleActivator active bundle as service in BundleContext
type BundleActivator interface {
	Start(ctx BundleContext)
	Stop(ctx BundleContext)
}

// BundleContext execute result
type BundleContext interface {
	GetServiceReference(serviceName string) (BundleServiceReference, error)
	GetBundles() ([]*Bundle, error)
	GetBundle(bundleName string) (*Bundle, error)
	InstallBundle(bundleName string) error
	UninstallBundle(bundleName string) error
	RegisterService(serviceName string, srv BundleService) error
	UnregisterService(serviceName string) error
	Stop() error
}
