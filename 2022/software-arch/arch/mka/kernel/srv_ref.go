package kernel

import "software-arch/arch/mka/bundle"

type DefaultBundleSrvRef struct {
	srvChan chan<- bundle.BundleServiceMessage
}

func NewDefaultBundleSrvRef(srvChan chan<- bundle.BundleServiceMessage) bundle.BundleServiceReference {
	return &DefaultBundleSrvRef{srvChan: srvChan}
}

func (ref *DefaultBundleSrvRef) Send(msg bundle.BundleServiceMessage) <-chan bundle.Message {
	resMsg := NewMessageWithResult(msg)
	ref.srvChan <- resMsg
	return resMsg.GetResultChan()
}
