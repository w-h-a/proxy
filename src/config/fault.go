package config

type Fault struct {
	Name   string
	Config RawConfig
}

type RawConfig map[string]interface{}
