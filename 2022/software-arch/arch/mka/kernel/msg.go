package kernel

import "software-arch/arch/mka/bundle"

type defaultMessage struct {
	functionName string
	message      interface{}
}

func NewDefaultMessage(functionName string, message interface{}) bundle.BundleServiceMessage {
	return &defaultMessage{
		functionName: functionName,
		message:      message,
	}
}

func (d *defaultMessage) GetFunctionName() string {
	return d.functionName
}

func (d *defaultMessage) GetMessage() interface{} {
	return d.message
}

type messageWithResult struct {
	bundle.BundleServiceMessage
	resChan chan bundle.Message
}

func NewMessageWithResult(msg bundle.BundleServiceMessage) *messageWithResult {
	return &messageWithResult{
		BundleServiceMessage: msg,
		resChan:              make(chan bundle.Message),
	}
}

func (msg *messageWithResult) GetResultChan() chan bundle.Message {
	return msg.resChan
}
