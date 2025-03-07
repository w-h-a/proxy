package http

import (
	"net/http"
	"net/http/httputil"

	"github.com/w-h-a/proxy/src/services/fault"
)

type Proxy struct {
	target string
	faults []fault.Fault
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	director := func(req *http.Request) {
		req = r
		req.URL.Host = p.target
	}

	for _, f := range p.faults {
		f.HandleEvent(fault.PRE_DISPATCH)
	}

	proxy := &httputil.ReverseProxy{Director: director}

	proxy.ServeHTTP(w, r)
}

func NewProxy(target string, faults []fault.Fault) *Proxy {
	return &Proxy{
		target: target,
		faults: faults,
	}
}
