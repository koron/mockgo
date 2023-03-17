package main

import "testing"

func TestMockFilename(t *testing.T) {
	forTest = false
	for i, tc := range []struct {
		typn string
		want string
	}{
		{"foo", "foo_mock.go"},
		{"FOO", "foo_mock.go"},
		{"foomock", "foo_mock.go"},
		{"FOOMOCK", "foo_mock.go"},
		{"foo_mock", "foo__mock.go"},
	} {
		got := mockFilename(tc.typn)
		if got != tc.want {
			t.Errorf("failed #%d %+v: got=%s", i, tc, got)
		}
	}
}
