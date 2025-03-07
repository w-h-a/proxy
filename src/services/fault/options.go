package fault

// type FaultOption func(o *FaultOptions)

type FaultOptions struct {
	Delay int `mapstructure:",omitempty"`
}

// func NewFaultOptions(opts ...FaultOption) FaultOptions {
// 	options := FaultOptions{}

// 	for _, fn := range opts {
// 		fn(&options)
// 	}

// 	return options
// }
