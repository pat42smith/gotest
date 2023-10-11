// Copyright 2022-2023 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"fmt"
	"testing"
)

func TestStubHelper(t *testing.T) {
	var sr StubReporter
	sr.Helper()
	sr.Expect(t, false, false, "")

	sr.failed = true
	sr.killed = true
	fmt.Fprintln(&sr.log, "boo")
	sr.Helper()
	sr.Expect(t, true, true, "boo\n")
}

func TestStubExpect(t *testing.T) {
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
		sr.Expect(&srt, f, k, l)
		if srt.failed != (bits>>3 != bits&7) {
			t.Errorf("Wrong answer for bits %o", bits)
		}
	}
}

func TestStubFail(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "")
	sr.Fail()
	sr.Expect(t, true, false, "")
}

func TestStubFailNow(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "")
	sr.FailNow()
	sr.Expect(t, true, true, "")
}

func TestStubLog(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "")
	sr.Log("one")
	sr.Expect(t, false, false, "one\n")
	sr.Log("two", "three\n")
	sr.Expect(t, false, false, "one\ntwo three\n\n")
	sr.Log()
	sr.Expect(t, false, false, "one\ntwo three\n\n\n")
}

func TestStubLogf(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "")
	sr.Logf("m%sy", "one")
	sr.Expect(t, false, false, "money\n")
	sr.Logf("")
	sr.Expect(t, false, false, "money\n\n")
	sr.Logf("%s\n", "two")
	sr.Expect(t, false, false, "money\n\ntwo\n")
}

func TestStubError(t *testing.T) {
	var sr StubReporter
	sr.Error("boo")
	sr.Expect(t, true, false, "boo\n")
}

func TestStubErrorf(t *testing.T) {
	var sr StubReporter
	sr.Errorf("boo")
	sr.Expect(t, true, false, "boo\n")
}

func TestStubFatal(t *testing.T) {
	var sr StubReporter
	sr.Fatal("boo")
	sr.Expect(t, true, true, "boo\n")
}

func TestStubFatalf(t *testing.T) {
	var sr StubReporter
	sr.Fatalf("boo")
	sr.Expect(t, true, true, "boo\n")
}

func TestStubReset(t *testing.T) {
	var sr StubReporter
	sr.Expect(t, false, false, "")
	sr.Fatal("boo")
	sr.Expect(t, true, true, "boo\n")
	sr.Reset()
	sr.Expect(t, false, false, "")
	sr.Reset()
	sr.Expect(t, false, false, "")
}

func TestStubMessages(t *testing.T) {
	var sr, x StubReporter
	sr.Expect(&x, true, true, "oops\n")
	x.Expect(t, true, true, `StubReporter marked not failed
StubReporter marked not killed
StubReporter log is ''; expected 'oops
'
`)

	sr.Reset()
	x.Reset()
	sr.Fatal("run!")
	sr.Expect(&x, false, false, "walk\n")
	x.Expect(t, true, true, `StubReporter marked failed
StubReporter marked killed
StubReporter log is 'run!
'; expected 'walk
'
`)
}
