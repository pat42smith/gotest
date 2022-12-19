// Copyright 2022 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"fmt"
	"testing"
)

func TestHelper(t *testing.T) {
	var sr StubReporter
	sr.Helper()
	// If it returned, we're fine.
}

func TestSRExpect(t *testing.T) {
	for bits := 0; bits < 0100; bits++ {
		var sr StubReporter
		if bits&1 != 0 {
			sr.failed = true
		}
		if bits&2 != 0 {
			sr.killed = true
		}
		if bits&4 != 0 {
			fmt.Fprintln(&sr.log, "one")
		} else {
			fmt.Fprintln(&sr.log, "zero")
		}

		f, k, l := false, false, "zero\n"
		if bits&010 != 0 {
			f = true
		}
		if bits&020 != 0 {
			k = true
		}
		if bits&040 != 0 {
			l = "one\n"
		}

		var srt StubReporter
		sr.Expect(&srt, f, k, l, "in TestSRExpect")
		if srt.failed != (bits>>3 != bits&7) {
			t.Errorf("Wrong answer for bits %o", bits)
		}
	}
}

func TestFail(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "", "before Fail")
	sr.Fail()
	sr.Expect(t, true, false, "", "after Fail")
}

func TestFailNow(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "", "before FailNow")
	sr.FailNow()
	sr.Expect(t, true, true, "", "after FailNow")
}

func TestLog(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "", "")
	sr.Log("one")
	sr.Expect(t, false, false, "one\n", "")
	sr.Log("two", "three\n")
	sr.Expect(t, false, false, "one\ntwo three\n\n", "")
	sr.Log()
	sr.Expect(t, false, false, "one\ntwo three\n\n\n", "")
}

func TestLogf(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "", "")
	sr.Logf("m%sy", "one")
	sr.Expect(t, false, false, "money\n", "")
	sr.Logf("")
	sr.Expect(t, false, false, "money\n\n", "")
	sr.Logf("%s\n", "two")
	sr.Expect(t, false, false, "money\n\ntwo\n", "")
}

func TestError(t *testing.T) {
	var sr StubReporter
	sr.Error("boo")
	sr.Expect(t, true, false, "boo\n", "")
}

func TestErrorf(t *testing.T) {
	var sr StubReporter
	sr.Errorf("boo")
	sr.Expect(t, true, false, "boo\n", "")
}

func TestFatal(t *testing.T) {
	var sr StubReporter
	sr.Fatal("boo")
	sr.Expect(t, true, true, "boo\n", "")
}

func TestFatalf(t *testing.T) {
	var sr StubReporter
	sr.Fatalf("boo")
	sr.Expect(t, true, true, "boo\n", "")
}

func TestReset(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "", "")
	sr.Fatal("boo")
	sr.Expect(t, true, true, "boo\n", "")
	sr.Reset()
	sr.Expect(t, false, false, "", "")
}
