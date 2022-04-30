package bundle

type BundleServiceReference interface {
	Send(BundleServiceMessage) <-chan Message
}

type BundleService interface {
	Recv(BundleServiceMessage) Message
}
