// Copyright 2022-2023 Patrick Smith
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
// be changed with the NotFatal wrapper.
package gotest

// Type Reporter is an interface satisfied by the testing.T, .B, and .F types.
//
// Reporter includes the methods involved in reporting the status of test cases.
// It also includes Helper, so that the helper functions defined in this package
// can properly mark themselves as helper functions.
//
// In future, more methods from the intersection of T, B, and F may be added.
// Be warned that this may break code containing types designed to implement Reporter;
// you create such types at your own risk.
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
		t.Fatal("Required condition failed")
	}
}

// Expect(t, a, b) is equivalent to Require(t, b == a), but with better messaging.
func Expect[T comparable](t Reporter, expected, actual T) {
	t.Helper()
	if actual != expected {
		t.Fatal("Expected", expected, "but actual value was", actual)
	}
}

// Function panics runs f and reports whether it panics.
//
// If f panics, panics returns true and the value passed to panic.
// Otherwise, panics returns (false, nil).
// If f calls panic(nil), panics will return (true, nil) in Go 1.20 and
// earlier. It will return true and a non-nil value in 1.21 and later.
func panics(f func()) (panicked bool, with any) {
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
// However, if f calls panic(nil), MustPanic will return a non-nil value
// in Go 1.21 and later versions.
func MustPanic(t Reporter, f func()) any {
	t.Helper()
	panicked, with := panics(f)
	if !panicked {
		t.Fatal("Expected panic did not occur")
	}
	return with
}

// NotFatal wraps a Reporter, and redirects fatal errors to non-terminating errors.
type NotFatal struct {
	Reporter
}

// NotFatal.FailNow marks a test failed, but does not not terminate the test;
// it is equivalent to nf.Fail().
func (nf NotFatal) FailNow() {
	nf.Fail()
}

// NotFatal.Fatal reports an error without terminating the test;
// it is equivalent to nf.Error(args...).
func (nf NotFatal) Fatal(args ...any) {
	nf.Error(args...)
}

// NotFatal.Fatal reports an error without terminating the test;
// it is equivalent to nf.Errorf(format, args...).
func (nf NotFatal) Fatalf(format string, args ...any) {
	nf.Errorf(format, args...)
}
