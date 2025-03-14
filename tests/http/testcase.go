package http

import "net/http"

type TestCase struct {
	When        string
	Faults      string
	Endpoint    string
	Query       string
	Then        string
	Status      int
	Header      http.Header
	Response    string
	DurationGTE float64
	DurationLTE float64
}
