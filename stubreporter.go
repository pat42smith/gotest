// Copyright 2022 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"fmt"
	"strings"
)

// Type *StubReporter is a simple implementation of the Reporter interface.
//
// It is intended to assist in testing test helper functions.
// The methods save the results of calls to be queried later,
// but do not do anything else.
type StubReporter struct {
	log            strings.Builder
	failed, killed bool
}

// Helper marks a function as a helper function.
//
// The StubReporter version of Helper does nothing.
func (sr *StubReporter) Helper() {}

// Fail marks a test as failed.
func (sr *StubReporter) Fail() {
	sr.failed = true
}

// Failed returns whether the test was marked failed.
func (sr *StubReporter) Failed() bool {
	return sr.failed
}

// FailNow marks a test as failed.
//
// The versions of FailNow in testing also terminate the test case, so they do not return.
// However, the FailNow, Fatal, and Fatalf methods of StubReporter do return.
// Careful attention should be paid to this point when using StubReporter
// to write test cases for helper functions.
func (sr *StubReporter) FailNow() {
	sr.Fail()
	sr.killed = true
}

// Killed returns whether FailNow was called.
func (sr *StubReporter) Killed() bool {
	return sr.killed
}

// Log formats its arguments as if by fmt.Println and records the resulting text.
func (sr *StubReporter) Log(args ...any) {
	_, e := fmt.Fprintln(&sr.log, args...)
	if e != nil {
		// Should be impossible
		panic(e)
	}
}

// Logf formats its arguments as if by fmt.Printf and records the resulting text.
func (sr *StubReporter) Logf(format string, args ...any) {
	oldLen := sr.log.Len()
	_, e := fmt.Fprintf(&sr.log, format, args...)
	if e != nil {
		// Should be impossible
		panic(e)
	}
	if sr.log.Len() == oldLen || !strings.HasSuffix(sr.log.String(), "\n") {
		sr.log.WriteByte('\n')
	}
}

// Logged returns the text recorded by Log and Logf.
func (sr *StubReporter) Logged() string {
	return sr.log.String()
}

// Error calls Fail and Log.
func (sr *StubReporter) Error(args ...any) {
	sr.Fail()
	sr.Log(args...)
}

// Errorf calls Fail and Logf.
func (sr *StubReporter) Errorf(format string, args ...any) {
	sr.Fail()
	sr.Logf(format, args...)
}

// Fatal calls FailNow and Log.
func (sr *StubReporter) Fatal(args ...any) {
	sr.FailNow()
	sr.Log(args...)
}

// Fatalf calls FailNow and Logf.
func (sr *StubReporter) Fatalf(format string, args ...any) {
	sr.FailNow()
	sr.Logf(format, args...)
}

// Expect verifies the status of the StubReporter.
//
// The failed, killed, and log parameters are compared to the StubReporter status.
// If they do not match, an error is reported to t; the message will include
// the when parameter.
func (sr *StubReporter) Expect(t Reporter, failed, killed bool, log, when string) {
	t.Helper()
	if sr.Failed() != failed {
		if !failed {
			t.Fatal("StubReporter marked failed", when)
		} else {
			t.Fatal("StubReporter marked not failed", when)
		}
	}
	if sr.Killed() != killed {
		if !killed {
			t.Fatal("StubReporter marked killed", when)
		} else {
			t.Fatal("StubReporter marked not killed", when)
		}
	}
	if actual := sr.Logged(); actual != log {
		t.Fatalf("%s StubReporter log is '%s'; expected '%s'", when, actual, log)
	}
}
