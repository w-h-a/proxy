package httpdelay

import (
	"context"
	"time"

	"github.com/w-h-a/proxy/src/services/fault"
)

type httpDelay struct {
	options fault.Options
}

func (f *httpDelay) Options() fault.Options {
	return f.options
}

func (f *httpDelay) HandleEvent(ctx context.Context, event fault.ProxyEvent) {
	switch event {
	case fault.PRE_DISPATCH:
		delay := time.Duration(f.options.Delay) * time.Second
		time.Sleep(delay)
	}
}

func NewFault(options fault.Options) fault.Fault {
	return &httpDelay{
		options: options,
	}
}
