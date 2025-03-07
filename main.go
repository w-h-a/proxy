package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"github.com/w-h-a/proxy/src"
	"github.com/w-h-a/proxy/src/config"
	"github.com/w-h-a/proxy/src/services/fault"
	httpdelay "github.com/w-h-a/proxy/src/services/fault/httpDelay"
)

func main() {
	// config
	config.NewConfig()

	// name
	name := fmt.Sprintf("%s.%s", config.Namespace(), config.Name())

	// log
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix(name),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)

	// traces

	// other dependencies
	httpClient := http.DefaultTransport

	faultManager := fault.NewManager()

	faultManager.Register(func(options fault.FaultOptions) (fault.Fault, error) {
		return httpdelay.NewFault(options), nil
	}, "httpdelay")

	// servers
	httpServer := src.AppFactory(httpClient, faultManager)

	// wait group and error chan
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)

	// start http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- httpServer.Start()
	}()

	// block
	err := <-errCh
	if err != nil {
		log.Errorf("failed to start server: %+v", err)
	}

	// graceful shutdown
	wait := make(chan struct{})

	go func() {
		defer close(wait)
		wg.Wait()
	}()

	select {
	case <-wait:
	case <-time.After(30 * time.Second):
	}

	log.Info("successfully stopped server")
}
