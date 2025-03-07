package config

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/w-h-a/pkg/telemetry/log"
)

var (
	instance *config
	once     sync.Once
)

type config struct {
	namespace           string
	name                string
	version             string
	httpAddress         string
	httpTargetScheme    string
	httpTargetNamespace string
	httpTargetName      string
	httpTargetPort      int
	faults              []Fault
}

func NewConfig() {
	once.Do(func() {
		instance = &config{
			namespace:           "dev",
			name:                "proxy",
			version:             "0.1.0-alpha.0",
			httpAddress:         ":0",
			httpTargetScheme:    "http",
			httpTargetNamespace: "dev",
			httpTargetName:      "localhost",
			httpTargetPort:      9090,
			faults:              []Fault{},
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

		httpTargetScheme := os.Getenv("HTTP_TARGET_SCHEME")
		if len(httpTargetScheme) > 0 {
			instance.httpTargetScheme = httpTargetScheme
		}

		httpTargetNamespace := os.Getenv("HTTP_TARGET_NAMESPACE")
		if len(httpTargetNamespace) > 0 {
			instance.httpTargetNamespace = httpTargetNamespace
		}

		httpTargetName := os.Getenv("HTTP_TARGET_NAME")
		if len(httpTargetName) > 0 {
			instance.httpTargetName = httpTargetName
		}

		httpTargetPort := os.Getenv("HTTP_TARGET_PORT")
		if len(httpTargetPort) > 0 {
			instance.httpTargetPort, _ = strconv.Atoi(httpTargetPort)
		}

		faultsJSON := os.Getenv("FAULTS")
		if len(faultsJSON) > 0 {
			var faults []Fault

			if err := json.Unmarshal([]byte(faultsJSON), &faults); err != nil {
				instance.faults = []Fault{}
			} else {
				instance.faults = faults
			}
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

func HttpTargetScheme() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.httpTargetScheme
}

func HttpTargetNamespace() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.httpTargetNamespace
}

func HttpTargetName() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.httpTargetName
}

func HttpTargetPort() int {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.httpTargetPort
}

func Faults() []Fault {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.faults
}

// For tests only
func Reset() {
	instance = &config{
		namespace:           "dev",
		name:                "proxy",
		version:             "0.1.0-alpha.0",
		httpAddress:         ":0",
		httpTargetScheme:    "http",
		httpTargetNamespace: "dev",
		httpTargetName:      "localhost",
		httpTargetPort:      9090,
		faults:              []Fault{},
	}

	once = sync.Once{}
}
