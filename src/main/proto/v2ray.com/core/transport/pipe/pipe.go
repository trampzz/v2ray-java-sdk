package pipe

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/signal/done"
)

// Option for creating new Pipes.
type Option func(*pipe)

func WithoutSizeLimit() Option {
	return func(p *pipe) {
		p.limit = -1
	}
}

func WithSizeLimit(limit int32) Option {
	return func(p *pipe) {
		p.limit = limit
	}
}

func OptionsFromContext(ctx context.Context) []Option {
	var opt []Option

	bp := core.BufferPolicyFromContext(ctx)
	if bp.PerConnection >= 0 {
		opt = append(opt, WithSizeLimit(bp.PerConnection))
	} else {
		opt = append(opt, WithoutSizeLimit())
	}

	return opt
}

// New creates a new Reader and Writer that connects to each other.
func New(opts ...Option) (*Reader, *Writer) {
	p := &pipe{
		limit:       -1,
		readSignal:  signal.NewNotifier(),
		writeSignal: signal.NewNotifier(),
		done:        done.New(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return &Reader{
			pipe: p,
		}, &Writer{
			pipe: p,
		}
}

type closeError interface {
	CloseError()
}

// CloseError invokes CloseError() method if the object is either Reader or Writer.
func CloseError(v interface{}) {
	if c, ok := v.(closeError); ok {
		c.CloseError()
	}
}
