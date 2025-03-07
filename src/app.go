package src

import (
	"fmt"
	"net/url"

	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/gorilla/mux"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/proxy/src/config"
	httphandlers "github.com/w-h-a/proxy/src/handlers/http"
	"github.com/w-h-a/proxy/src/services/fault"
)

func AppFactory(manager *fault.Manager) serverv2.Server {
	// faults
	faults := []fault.Fault{}

	for _, f := range config.Faults() {
		factory, ok := manager.Lookup(f.Name)
		if !ok {
			log.Fatalf("failed to lookup fault %s", f.Name)
		}

		options := fault.FaultOptions{}

		if err := mapstructure.Decode(f.Config, &options); err != nil {
			log.Fatalf("failed to apply fault configuration for %s: %v", f.Name, err)
		}

		faultService, err := factory(options)
		if err != nil {
			log.Fatalf("invoking the fault factory for %s resulted in an error: %v", f.Name, err)
		}

		faults = append(faults, faultService)
	}

	// base server options
	opts := []serverv2.ServerOption{
		serverv2.ServerWithNamespace(config.Namespace()),
		serverv2.ServerWithName(config.Name()),
		serverv2.ServerWithVersion(config.Version()),
	}

	// create http server
	router := mux.NewRouter()

	targetURL, err := url.Parse(fmt.Sprintf("http://%s", config.HttpTarget()))
	if err != nil {
		log.Fatalf("failed to parse target url %s: %v", config.HttpTarget(), err)
	}

	proxy := httphandlers.NewProxy(targetURL, faults)

	router.PathPrefix("/").HandlerFunc(proxy.Serve)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	httpServer.Handle(router)

	return httpServer
}
