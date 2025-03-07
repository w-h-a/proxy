package http

type TestCase struct {
	When        string
	Endpoint    string
	Query       string
	Response    string
	Faults      string
	Then        string
	DurationGTE float64
	DurationLTE float64
}
