/*
Package mockrt provides mock runtime for mockgo
*/
package mockrt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Call defines pair of request and response parameters for a method call.
type Call struct {
	Request  interface{}
	Response interface{}
}

// Sequence is a checker of sequence of method calls
type Sequence struct {
	t     *testing.T
	calls []Call
	opts  []cmp.Option
	index int
}

// NewSequence creates a sequence of calls.
// This is called by test codes.
func NewSequence(t *testing.T, calls ...Call) *Sequence {
	return &Sequence{
		t:     t,
		calls: calls,
	}
}

// AddCall adds call data.
// This is called by test codes.
func (s *Sequence) AddCall(calls ...Call) *Sequence {
	s.calls = append(s.calls, calls...)
	return s
}

// WithOption updates compare option.
// This is called by test codes.
func (s *Sequence) WithOption(opts ...cmp.Option) *Sequence {
	s.opts = opts
	return s
}

// Call checks call parameters (request) and return response.
// This is called by mock code.
func (s *Sequence) Call(name string, req interface{}) interface{} {
	s.t.Helper()
	if s.index >= len(s.calls) {
		s.t.Fatalf("no calls at #%d for %s\nreq=%+v", s.index, name, req)
	}
	c := s.calls[s.index]
	if d := cmp.Diff(c.Request, req, s.opts...); d != "" {
		s.t.Fatalf("call for %s (#%d) has unexpected arguments: -want +got\n%s", name, s.index, d)
	}
	s.index++
	return c.Response
}

// IsEnd checks sequence has end or not.
// This is called by test code.
func (s *Sequence) IsEnd() {
	s.t.Helper()
	if s.index < len(s.calls) {
		s.t.Fatalf("there are non-proceeded calles: %+v", s.calls[s.index:])
	}
}
