package http

import (
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	testCases := []TestCase{
		{
			When:        "when: there are no faults, we call the customer endpoint with a query, and the backend responds with 'I am the backend' and 200",
			Faults:      `[]`,
			Endpoint:    "/customer",
			Query:       "?id=30",
			Then:        "then: the proxy responds with 'I am the backend' and 200 with normal latency",
			Status:      200,
			Response:    "I am the backend",
			DurationGTE: 0,
			DurationLTE: 2,
		},
		{
			When:        "when: there is an http delay for GET /customer, we call the customer endpoint with a query, and the backend responds with 'I am the backend' and 200",
			Faults:      `[{"name":"httpdelay","config":{"delay":2,"rules":[{"endpoint":"/cust.*","httpmethod":"(?i)gEt"}]}}]`,
			Endpoint:    "/customer",
			Query:       "?id=30",
			Then:        "then: the proxy responds with 'I am the backend' and 200 with 2s latency",
			Status:      200,
			Response:    "I am the backend",
			DurationGTE: 2,
			DurationLTE: 4,
		},
		{
			When:        "when: there is an http delay for GET /customer, we call the customer endpoint with a query, and we tamper with the status code, headers, and body",
			Faults:      `[{"name":"httptamper","config":{"status":500,"headers":{"foo-bar":"baz"},"body":"I screwed up"}},{"name":"httpdelay","config":{"delay":1}}]`,
			Endpoint:    "/customer",
			Query:       "?id=30",
			Then:        "then: the proxy responds with 'I screwed up', 'foo-bar' header, and 500 with 1s latency",
			Status:      500,
			Header:      http.Header{"foo-bar": []string{"baz"}},
			Response:    "I screwed up",
			DurationGTE: 1,
			DurationLTE: 3,
		},
		{
			When:        "when: we call the customer endpoint with a query, but we tamper with the status code and body of /bar.*",
			Faults:      `[{"name":"httptamper","config":{"status":500,"body":"I screwed up","rules":[{"endpoint":"/bar.*"}]}}]`,
			Endpoint:    "/customer",
			Query:       "?id=30",
			Then:        "then: the proxy responds with 'I am the backend' and 200",
			Status:      200,
			Response:    "I am the backend",
			DurationGTE: 0,
			DurationLTE: 2,
		},
	}

	RunTestCases(t, testCases)
}
