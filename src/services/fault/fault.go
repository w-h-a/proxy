package fault

import "context"

type ProxyEvent int

const (
	PRE_DISPATCH ProxyEvent = iota
	POST_DISPATCH
)

type Fault interface {
	Options() Options
	HandleEvent(ctx context.Context, event ProxyEvent)
}
