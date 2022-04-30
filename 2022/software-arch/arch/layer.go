package arch

import (
	"context"
	"fmt"
	"io"
)

func NewTopLayer(next MiddleExecutor) LayerExecutor {
	return &TopLayer{nextLayer: next}
}

func NewMiddleLayer(next LowExecutor) MiddleExecutor {
	return &MiddleLayer{nextLayer: next}
}

func NewLowLayer(ioW io.Writer) LowExecutor {
	return &LowLayer{ioW: ioW}
}

type LayerExecutor interface {
	Do(ctx context.Context) error
}

type MiddleExecutor interface {
	DoMiddle(ctx context.Context) error
}

type LowExecutor interface {
	DoLow(ctx context.Context) error
}

type TopLayer struct {
	nextLayer MiddleExecutor
}

func (t *TopLayer) Do(ctx context.Context) error {
	c, ok := ctx.Value(GLabelContent).(string)
	if !ok {
		return fmt.Errorf("%s not exists", GLabelContent)
	}
	return t.nextLayer.DoMiddle(context.WithValue(ctx, GLabelContent, fmt.Sprintf("#TopLayer[%s]", c)))
}

type MiddleLayer struct {
	nextLayer LowExecutor
}

func (t *MiddleLayer) DoMiddle(ctx context.Context) error {
	c, ok := ctx.Value(GLabelContent).(string)
	if !ok {
		return fmt.Errorf("%s not exists", GLabelContent)
	}
	return t.nextLayer.DoLow(context.WithValue(ctx, GLabelContent, fmt.Sprintf("#MiddleLayer[%s]", c)))
}

type LowLayer struct {
	ioW io.Writer
}

func (t *LowLayer) DoLow(ctx context.Context) error {
	c, ok := ctx.Value(GLabelContent).(string)
	if !ok {
		return fmt.Errorf("%s not exists", GLabelContent)
	}
	n, err := fmt.Fprintf(t.ioW, "#LowLayer[%s]", c)
	if n != (len(c) + 11) {
		return fmt.Errorf("not all wrotten")
	}
	return err
}
