package config

import (
	"os"
	"sync"

	"github.com/w-h-a/pkg/telemetry/log"
)

var (
	instance *config
	once     sync.Once
)

type config struct {
	namespace   string
	name        string
	version     string
	httpAddress string
	httpTarget  string
}

func NewConfig() {
	once.Do(func() {
		instance = &config{
			namespace:   "dev",
			name:        "proxy",
			version:     "0.1.0-alpha.0",
			httpAddress: ":0",
			httpTarget:  "localhost:9090",
		}

		namespace := os.Getenv("NAMESPACE")
		if len(namespace) > 0 {
			instance.namespace = namespace
		}

		name := os.Getenv("NAME")
		if len(name) > 0 {
			instance.name = name
		}

		version := os.Getenv("VERSION")
		if len(version) > 0 {
			instance.version = version
		}

		httpAddress := os.Getenv("HTTP_ADDRESS")
		if len(httpAddress) > 0 {
			instance.httpAddress = httpAddress
		}

		httpTarget := os.Getenv("HTTP_TARGET")
		if len(httpTarget) > 0 {
			instance.httpTarget = httpTarget
		}
	})
}

func Namespace() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.namespace
}

func Name() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.name
}

func Version() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.version
}

func HttpAddress() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.httpAddress
}

func HttpTarget() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.httpTarget
}
