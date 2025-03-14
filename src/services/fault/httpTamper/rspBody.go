package httptamper

import "io"

type rspBody struct {
	body   []byte
	closed bool
}

func (r *rspBody) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.EOF
	}

	copy(p, r.body)

	r.closed = true

	return len(r.body), nil
}

func (r *rspBody) Close() error {
	r.closed = true
	return nil
}
