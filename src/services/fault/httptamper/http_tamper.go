package httptamper

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/w-h-a/proxy/src/services/fault"
)

type httpTamper struct {
	options fault.Options
}

func (f *httpTamper) Options() fault.Options {
	return f.options
}

func (f *httpTamper) HandleEvent(ctx context.Context, event fault.ProxyEvent) {
	if !fault.ShouldApply(ctx, f.options.Rules) {
		return
	}

	switch event {
	case fault.POST_DISPATCH:
		rsp, _ := fault.HttpResponseFromCtx(ctx)

		for k, v := range f.options.Headers {
			rsp.Header.Set(k, v)
		}

		if f.options.Status != 0 {
			rsp.StatusCode = f.options.Status
			rsp.Status = http.StatusText(f.options.Status)
		}

		if f.options.Body != "" {
			if rsp.Body != nil {
				_, _ = io.ReadAll(rsp.Body)
				rsp.Body.Close()
			}

			bs := []byte(f.options.Body)

			body := io.NopCloser(bytes.NewBuffer(bs))

			rsp.Body = body
			rsp.ContentLength = int64(len(bs))
		}
	}
}

func NewFault(options fault.Options) fault.Fault {
	return &httpTamper{
		options: options,
	}
}
