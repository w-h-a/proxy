package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
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
	faults      []Fault
}

func NewConfig() {
	once.Do(func() {
		instance = &config{
			namespace:   "dev",
			name:        "proxy",
			version:     "0.1.0-alpha.0",
			httpAddress: ":0",
			httpTarget:  "localhost:9090",
			faults:      []Fault{},
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

		faultsJSON := os.Getenv("FAULTS")
		if len(faultsJSON) > 0 {
			var faults []Fault

			if err := json.Unmarshal([]byte(faultsJSON), &faults); err != nil {
				log.Fatalf("failed to unmarshal %s faults JSON config: %v", faultsJSON, err)
			}

			instance.faults = faults
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

func Faults() []Fault {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.faults
}
