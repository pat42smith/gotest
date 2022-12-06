// Copyright 2022 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"strings"
	"testing"
)

func TestRunTrue(t *testing.T) {
	var st StubReporter
	RunCommand(&st, "/bin/true")
	st.Expect(t, false, false, "", "")
}

func TestRunFalse(t *testing.T) {
	var st StubReporter
	RunCommand(&st, "/bin/false")
	st.Expect(t, true, true, "/bin/false: exit status 1\n", "")
}

func TestRunShExit0(t *testing.T) {
	var st StubReporter
	RunCommand(&st, "/bin/sh", "-c", "exit 0")
	st.Expect(t, false, false, "", "")
}

func TestRunShExit5(t *testing.T) {
	var st StubReporter
	RunCommand(&st, "/bin/sh", "-c", "exit 5")
	st.Expect(t, true, true, "/bin/sh: exit status 5\n", "")
}

func TestRunNoneSuch(t *testing.T) {
	var st StubReporter
	cmd := "/bin/thereisnocommandwiththisname-really"
	RunCommand(&st, cmd)
	Require(t, st.Failed() && st.Killed())
	Require(t, st.Logged() != "")
}

func TestOutput(t *testing.T) {
	var st StubReporter
	RunCommand(&st, "/bin/sh", "-c", "echo something")
	st.Expect(t, true, true, "/bin/sh: unexpected output:\nsomething\n", "")
}

func TestErrorOutput(t *testing.T) {
	var st StubReporter
	RunCommand(&st, "/bin/sh", "-c", "echo whatever >&2")
	st.Expect(t, true, true, "/bin/sh: unexpected output:\nwhatever\n", "")
}

func TestLong(t *testing.T) {
	r := strings.Repeat("x", 1200)
	input := "a" + r + "b"
	begin := "a" + r[:149]
	end := r[:148] + "b\n"
	want := begin + " ... " + end

	var st StubReporter
	RunCommand(&st, "/bin/sh", "-c", "echo "+input)
	st.Expect(t, true, true, "/bin/sh: unexpected output:\n"+want, "")
}

func TestDoubleFail(t *testing.T) {
	var st StubReporter
	RunCommand(&st, "/bin/sh", "-c", "echo failing; exit 13")
	st.Expect(t, true, true, "/bin/sh: unexpected output:\nfailing\n/bin/sh: exit status 13\n", "")
}
