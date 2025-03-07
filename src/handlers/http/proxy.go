package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/proxy/src/services/fault"
)

type Proxy struct {
	targetNamespace string
	targetName      string
	targetPort      int
	hopHeaders      []string
	faults          []fault.Fault
	client          http.RoundTripper
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// req := r.Clone(ctx)

	host := fmt.Sprintf("%s:%d", p.targetName, p.targetPort)

	// // TODO: make configurable
	// req.URL.Scheme = "http"
	// req.URL.Host = host
	// req.Host = host
	// req.Close = false

	// for k, h := range r.Header {
	// 	for _, v := range h {
	// 		req.Header.Add(k, v)
	// 	}
	// }

	// for _, h := range p.hopHeaders {
	// 	if req.Header.Get(h) != "" {
	// 		req.Header.Del(h)
	// 	}
	// }

	targetURL, _ := url.Parse("http://" + host)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	for _, f := range p.faults {
		f.HandleEvent(fault.PRE_DISPATCH)
	}

	log.Infof("AFTER %+v", r)

	proxy.ServeHTTP(w, r)

	// rsp, err := p.client.RoundTrip(req)
	// if err != nil {
	// 	log.Errorf("error: %+v", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// defer rsp.Body.Close()

	// for _, f := range p.faults {
	// 	f.HandleEvent(fault.POST_DISPATCH)
	// }

	// for _, h := range p.hopHeaders {
	// 	if rsp.Header.Get(h) != "" {
	// 		rsp.Header.Del(h)
	// 	}
	// }

	// rspHeader := w.Header()

	// for k, h := range rsp.Header {
	// 	for _, v := range h {
	// 		rspHeader.Add(k, v)
	// 	}
	// }

	// w.WriteHeader(rsp.StatusCode)
	// io.Copy(w, rsp.Body)
}

func NewProxy(targetNamespace string, targetName string, targetPort int, faults []fault.Fault, client http.RoundTripper) *Proxy {
	return &Proxy{
		targetNamespace: targetNamespace,
		targetName:      targetName,
		targetPort:      targetPort,
		hopHeaders:      []string{"connection", "keep-alive", "proxy-connection"},
		faults:          faults,
		client:          client,
	}
}
