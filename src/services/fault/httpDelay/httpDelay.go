package httpdelay

import (
	"time"

	"github.com/w-h-a/proxy/src/services/fault"
)

type httpDelay struct {
	options fault.FaultOptions
}

func (f *httpDelay) Options() fault.FaultOptions {
	return f.options
}

func (f *httpDelay) HandleEvent(event fault.ProxyEvent) {
	switch event {
	case fault.PRE_DISPATCH:
		delay := time.Duration(f.options.Delay) * time.Second
		time.Sleep(delay)
	}
}

func NewFault(options fault.FaultOptions) fault.Fault {
	return &httpDelay{
		options: options,
	}
}
