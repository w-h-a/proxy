package fault

type Options struct {
	Delay   int
	Status  int
	Headers map[string]string
	Body    string
}
