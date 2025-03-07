package fault

type ProxyEvent int

const (
	PRE_DISPATCH ProxyEvent = iota
	POST_DISPATCH
)

type Fault interface {
	Options() FaultOptions
	HandleEvent(event ProxyEvent)
}
