package fault

import (
	"context"
	"net/http"
)

type httpRequestKey struct{}

func CtxWithHttpRequest(ctx context.Context, req *http.Request) context.Context {
	ctx = context.WithValue(ctx, httpRequestKey{}, req)
	return ctx
}

func HttpRequestFromCtx(ctx context.Context) (*http.Request, bool) {
	r, ok := ctx.Value(httpRequestKey{}).(*http.Request)
	return r, ok
}

type httpResponseKey struct{}

func CtxWithHttpResponse(ctx context.Context, rsp *http.Response) context.Context {
	ctx = context.WithValue(ctx, httpResponseKey{}, rsp)
	return ctx
}

func HttpResponseFromCtx(ctx context.Context) (*http.Response, bool) {
	r, ok := ctx.Value(httpResponseKey{}).(*http.Response)
	return r, ok
}
