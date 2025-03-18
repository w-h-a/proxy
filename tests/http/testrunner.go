package http

import (
	"fmt"
	"io"
	gohttp "net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"github.com/w-h-a/proxy/src"
	httpclient "github.com/w-h-a/proxy/src/clients/http"
	"github.com/w-h-a/proxy/src/config"
)

func RunTestCases(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		// backend
		backend := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
			require.Equal(t, "", r.Header.Get("connection"))
			require.Equal(t, "", r.Header.Get("keep-alive"))

			if r.URL.Path == testCase.Endpoint {
				w.Write([]byte("I am the backend"))
			} else {
				gohttp.NotFound(w, r)
			}
		}))

		backendURL, err := url.Parse(backend.URL)
		require.NoError(t, err)

		target := backendURL.Host
		split := strings.Split(target, ":")

		// env vars
		os.Setenv("HTTP_TARGET_SCHEME", "http")
		os.Setenv("HTTP_TARGET_NAMESPACE", "test")
		os.Setenv("HTTP_TARGET_NAME", split[0])
		os.Setenv("HTTP_TARGET_PORT", split[1])
		os.Setenv("FAULTS", testCase.Faults)

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

		// clients
		httpClient := httpclient.NewHttpClient()

		// servers
		httpServer := src.AppFactory(httpClient)

		// tests
		t.Run(testCase.When, func(t *testing.T) {
			err = httpServer.Run()
			require.NoError(t, err)

			now := time.Now()

			rsp, err := gohttp.Get(fmt.Sprintf("http://%s%s%s", httpServer.Options().Address, testCase.Endpoint, testCase.Query))
			require.NoError(t, err)

			bs, err := io.ReadAll(rsp.Body)
			require.NoError(t, err)
			defer rsp.Body.Close()

			duration := time.Since(now).Seconds()

			t.Log(testCase.Then)

			require.GreaterOrEqual(t, duration, testCase.DurationGTE)
			require.LessOrEqual(t, duration, testCase.DurationLTE)

			require.Equal(t, testCase.Response, string(bs))

			require.Equal(t, testCase.Status, rsp.StatusCode)

			for k, v := range testCase.Header {
				require.Equal(t, v[0], rsp.Header.Get(k))
			}

			t.Cleanup(func() {
				err = httpServer.Stop()
				require.NoError(t, err)

				config.Reset()

				backend.Close()
			})
		})
	}
}
