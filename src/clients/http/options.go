package http

type Option func(*Options)

type Options struct {
}

func NewHttpClientOptions(opts ...Option) Options {
	options := Options{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
