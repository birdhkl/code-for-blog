package arch

type LabelType string

const (
	_ int = iota
	EventOnInput
	EventOnScrollAdded
	EventOnTextChanged
	EventOnScrollPicked
	EventOnScrollDel
)

const (
	GLabelEventType LabelType = "EventType"
	GLabelContent   LabelType = "Content"
)
