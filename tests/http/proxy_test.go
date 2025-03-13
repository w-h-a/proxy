package http

import (
	"testing"
)

func TestProxy(t *testing.T) {
	testCases := []TestCase{
		{
			When:        "when: there are no faults, we call the customer endpoint with a query, and the backend responds with 'I am the backend'",
			Faults:      `[]`,
			Endpoint:    "/customer",
			Query:       "?id=30",
			Response:    "I am the backend",
			Then:        "then: the proxy responds with 'I am the backend' with normal latency",
			DurationGTE: 0,
			DurationLTE: 2,
		},
		{
			When:        "when: there is an http delay, we call the customer endpoint with a query, and the backend responds with 'I am the backend'",
			Faults:      `[{"name":"httpdelay","config":{"delay":2}}]`,
			Endpoint:    "/customer",
			Query:       "?id=30",
			Response:    "I am the backend",
			Then:        "then: the proxy responds with 'I am the backend' with 2s latency",
			DurationGTE: 2,
			DurationLTE: 4,
		},
	}

	RunTestCases(t, testCases)
}
