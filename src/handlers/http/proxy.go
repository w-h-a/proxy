package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/w-h-a/proxy/src/services/fault"
)

type Proxy struct {
	target *url.URL
	faults []fault.Fault
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	for _, f := range p.faults {
		f.HandleEvent(fault.PRE_DISPATCH)
	}

	proxy := httputil.NewSingleHostReverseProxy(p.target)

	r.URL.Scheme = p.target.Scheme
	r.URL.Host = p.target.Host
	r.Host = p.target.Host

	proxy.ServeHTTP(w, r)
}

func NewProxy(target *url.URL, faults []fault.Fault) *Proxy {
	return &Proxy{
		target: target,
		faults: faults,
	}
}
