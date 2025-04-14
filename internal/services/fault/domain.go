package fault

type Rule struct {
	Endpoint   string
	HttpMethod string
	Percentage float64
}
