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
	Parameter interface{}
	Result    interface{}
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

// NewQ is an alias for NewSequence, creates a sequence of calls.
// This is called by test codes.
func NewQ(t *testing.T, calls ...Call) *Sequence {
	return NewSequence(t, calls...)
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

// Call checks call parameter and returns result.
// This is called by mock code.
func (s *Sequence) Call(name string, param interface{}) interface{} {
	s.t.Helper()
	if s.index >= len(s.calls) {
		s.t.Fatalf("no calls at #%d for %s\nparam=%+v", s.index, name, param)
	}
	c := s.calls[s.index]
	if d := cmp.Diff(c.Parameter, param, s.opts...); d != "" {
		s.t.Fatalf("call for %s (#%d) has unexpected arguments: -want +got\n%s", name, s.index, d)
	}
	s.index++
	return c.Result
}

// T returns *testing.T.
// This is called by mock code.
func (s *Sequence) T() *testing.T {
	return s.t
}

// IsEnd checks sequence has end or not.
// This is called by test code.
func (s *Sequence) IsEnd() {
	s.t.Helper()
	if s.index < len(s.calls) {
		s.t.Fatalf("there are non-proceeded calles: %+v", s.calls[s.index:])
	}
}
