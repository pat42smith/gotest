// Copyright 2022 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

// Package gotest provides support for writing test cases in Go.
//
// The intent is to make it easier to write test cases, even if this results
// in a lack of detail when reporting test failures. The hope is that if more
// detail is required to diagnose a failure, the test case can then be modified
// to report that detail.
//
// By default, these functions use FailNow, Fatal, or Fatalf to report errors,
// so they terminate the running test on the first error. This behavior may
// be changed with the ErrorsNotFatal wrapper.
package gotest

// Type Reporter is an interface satisfied by the testing.T, .B, and .F types.
//
// Reporter includes the methods involved in reporting the status of test cases.
// It also includes Helper, so that the helper functions defined in this package
// can properly mark themselves as helper functions.
//
// In future, more methods from the intersection of T, B, and F may be added.
// Be warned that code using Reporter in unusual ways may break if this happens.
type Reporter interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
}

// Require fails and terminates the running test if the condition is false.
func Require(t Reporter, condition bool) {
	t.Helper()
	if !condition {
		t.Fatal("Test requirement failed")
	}
}

// Expect(t, a, b) is equivalent to Require(t, b == a), but with better messaging.
func Expect[T comparable](t Reporter, expected, actual T) {
	t.Helper()
	if actual != expected {
		t.Fatal("Expected", expected, "but actual value was", actual)
	}
}

// NilError fails and terminates the test if passed a non-nil error.
func NilError(t Reporter, e error) {
	t.Helper()
	if e != nil {
		t.Fatal(e)
	}
}

// Panics runs f and reports whether it panics.
//
// If f panics, Panics returns true and the value passed to panic.
// Otherwise, Panics returns (false, nil).
// If f calls panic(nil), Panics will correctly return (true, nil).
func Panics(f func()) (panicked bool, with any) {
	defer func() {
		if panicked {
			with = recover()
		}
	}()
	panicked = true
	f()
	return false, nil
}

// MustPanic runs f and verifies that it panics.
//
// If f does not panic, MustPanic terminates the running test with an error.
// If f does panic, MustPanic returns the value that was passed to panic.
func MustPanic(t Reporter, f func()) any {
	t.Helper()
	panicked, with := Panics(f)
	if !panicked {
		t.Fatal("Expected panic did not occur")
	}
	return with
}

// ErrorsNotFatal wraps a Reporter, and redirects fatal errors to non-terminating errors.
type ErrorsNotFatal struct {
	Reporter
}

// ErrorsNotFatal.FailNow marks a test failed, but does not not terminate it.
func (enf ErrorsNotFatal) FailNow() {
	enf.Fail()
}

// ErrorsNotFatal.Fatal reports an error without terminating the test.
func (enf ErrorsNotFatal) Fatal(args ...any) {
	enf.Error(args...)
}

// ErrorsNotFatal.Fatal reports an error without terminating the test.
func (enf ErrorsNotFatal) Fatalf(format string, args ...any) {
	enf.Errorf(format, args...)
}
