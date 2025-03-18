package http

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/proxy/src/services/fault"
)

type Proxy struct {
	targetScheme    string
	targetNamespace string
	targetName      string
	targetPort      int
	hopHeaders      []string
	faults          []fault.Fault
	client          http.RoundTripper
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := r.Clone(ctx)

	req.URL.Scheme = p.targetScheme

	host := fmt.Sprintf("%s:%d", p.targetName, p.targetPort)
	req.URL.Host = host
	req.Host = host

	req.Close = false

	for k, h := range r.Header {
		for _, v := range h {
			req.Header.Add(k, v)
		}
	}

	for _, h := range p.hopHeaders {
		if req.Header.Get(h) != "" {
			req.Header.Del(h)
		}
	}

	for _, f := range p.faults {
		ctx := fault.CtxWithHttpRequest(context.Background(), req)

		f.HandleEvent(ctx, fault.PRE_DISPATCH)
	}

	rsp, err := p.client.RoundTrip(req)
	if err != nil {
		log.Errorf("error: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rsp.Body.Close()

	for _, f := range p.faults {
		ctx := fault.CtxWithHttpRequest(context.Background(), req)
		ctx = fault.CtxWithHttpResponse(ctx, rsp)

		f.HandleEvent(ctx, fault.POST_DISPATCH)
	}

	for _, h := range p.hopHeaders {
		if rsp.Header.Get(h) != "" {
			rsp.Header.Del(h)
		}
	}

	for k, h := range rsp.Header {
		if k != "Content-Length" {
			for _, v := range h {
				w.Header().Add(k, v)
			}
		}
	}

	w.WriteHeader(rsp.StatusCode)

	if rsp.Body != nil {
		bs, _ := io.ReadAll(rsp.Body)
		w.Write(bs)
	}
}

func NewProxy(
	targetScheme string,
	targetNamespace string,
	targetName string,
	targetPort int,
	faults []fault.Fault,
	client http.RoundTripper,
) *Proxy {
	return &Proxy{
		targetScheme:    targetScheme,
		targetNamespace: targetNamespace,
		targetName:      targetName,
		targetPort:      targetPort,
		hopHeaders: []string{
			"connection",
			"keep-alive",
			"proxy-connection",
		},
		faults: faults,
		client: client,
	}
}
