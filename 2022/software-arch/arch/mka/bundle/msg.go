package bundle

type Message interface {
	Error() error
	Value() interface{}
}

type BundleServiceMessage interface {
	GetFunctionName() string
	GetMessage() interface{}
}
