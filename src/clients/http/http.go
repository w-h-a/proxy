package http

import (
	gohttp "net/http"
)

type HttpClient struct {
	options      Options
	roundtripper gohttp.RoundTripper
}

func (c *HttpClient) Options() Options {
	return c.options
}

func (c *HttpClient) RoundTrip(req *gohttp.Request) (*gohttp.Response, error) {
	return c.roundtripper.RoundTrip(req)
}

func NewHttpClient(opts ...Option) gohttp.RoundTripper {
	options := NewHttpClientOptions(opts...)

	return &HttpClient{
		options:      options,
		roundtripper: gohttp.DefaultTransport,
	}
}
