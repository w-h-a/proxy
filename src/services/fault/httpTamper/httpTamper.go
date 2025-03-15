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
	switch event {
	case fault.POST_DISPATCH:
		var req *http.Request
		if r, ok := fault.HttpRequestFromCtx(ctx); !ok {
			return
		} else {
			req = r
		}

		var rsp *http.Response
		if r, ok := fault.HttpResponseFromCtx(ctx); !ok {
			return
		} else {
			rsp = r
		}

		if f.options.Body != "" {
			body := io.NopCloser(bytes.NewReader([]byte(f.options.Body)))

			r := &http.Response{
				Request:       req,
				Header:        rsp.Header,
				Close:         rsp.Close,
				ContentLength: rsp.ContentLength,
				TLS:           rsp.TLS,
				Status:        rsp.Status,
				StatusCode:    rsp.StatusCode,
				Body:          body,
			}

			*rsp = *r
		}

		for k, v := range f.options.Headers {
			rsp.Header.Set(k, v)
		}

		if f.options.Status != 0 {
			rsp.StatusCode = f.options.Status
			rsp.Status = http.StatusText(f.options.Status)
		}
	}
}

func NewFault(options fault.Options) fault.Fault {
	return &httpTamper{
		options: options,
	}
}
