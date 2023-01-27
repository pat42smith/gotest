// Copyright 2023 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCmdDefaults(t *testing.T) {
	var st StubReporter
	Command("/bin/true").Run(&st, "")
	st.Expect(t, false, false, "", "")

	st.Reset()
	Command("/bin/false").Run(&st, "")
	st.Expect(t, true, true, `non-zero exit code
command: /bin/false
no input
no output
no error output
exit code: 1
`, "")

	st.Reset()
	Command("/bin/printf", "99").Run(&st, "")
	st.Expect(t, true, true, `unexpected output
command: /bin/printf 99
no input
output:
99
no error output
exit code: 0
`, "")

	st.Reset()
	Command("/bin/sh", "-c", "echo 87 >&2").Run(&st, "")
	st.Expect(t, true, true, `unexpected error output
command: /bin/sh -c echo 87 >&2
no input
no output
error output:
87
exit code: 0
`, "")

	st.Reset()
	Command("/bin/sh", "-c", "echo 87 >&2; exit 3").Run(&st, "")
	st.Expect(t, true, true, `unexpected error output
command: /bin/sh -c echo 87 >&2; exit 3
no input
no output
error output:
87
exit code: 3
`, "")

	st.Reset()
	Command("/bin/sh", "-c", "echo 99; echo 87 >&2; exit 3").Run(&st, "")
	st.Expect(t, true, true, `unexpected output
unexpected error output
command: /bin/sh -c echo 99; echo 87 >&2; exit 3
no input
output:
99
error output:
87
exit code: 3
`, "")

	st.Reset()
	Command("/bin/sh", "-c", "echo 99; exit 3").Run(&st, "")
	st.Expect(t, true, true, `unexpected output
command: /bin/sh -c echo 99; exit 3
no input
output:
99
no error output
exit code: 3
`, "")
}

func TestCmdOutput(t *testing.T) {
	var st StubReporter
	c := Command("/bin/sh", "-c", "read x; echo a $x b")
	c.WantStdout("a seven b\n")

	c.Run(&st, "seven\n")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "eight\n")
	st.Expect(t, true, true, `incorrect output
command: /bin/sh -c read x; echo a $x b
input:
eight
output:
a eight b
no error output
exit code: 0
`, "")

	st.Reset()
	c.CheckStdout(func(actual string) bool {
		return strings.Contains(actual, "green")
	})
	c.Run(&st, "evergreen")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "purple")
	st.Expect(t, true, true, `incorrect output
command: /bin/sh -c read x; echo a $x b
input:
purple
output:
a purple b
no error output
exit code: 0
`, "")

	st.Reset()
	c.CheckStdout(nil)
	c.Run(&st, "whatever")
	st.Expect(t, true, true, `unexpected output
command: /bin/sh -c read x; echo a $x b
input:
whatever
output:
a whatever b
no error output
exit code: 0
`, "")
}

func TestCmdError(t *testing.T) {
	var st StubReporter
	c := Command("/bin/sh", "-c", "read x; if [ \"$x\" != nothing ]; then echo $x >&2; exit 99; fi")
	c.WantStderr("fever\n")

	c.Run(&st, "fever")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "chill")
	st.Expect(t, true, true, `incorrect error output
command: /bin/sh -c read x; if [ "$x" != nothing ]; then echo $x >&2; exit 99; fi
input:
chill
no output
error output:
chill
exit code: 99
`, "")

	st.Reset()
	c.CheckStderr(func(actual string) bool {
		return strings.Contains(actual, "tropical")
	})
	c.Run(&st, "blueberries are not a tropical fruit")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "apples grow in England")
	st.Expect(t, true, true, `incorrect error output
command: /bin/sh -c read x; if [ "$x" != nothing ]; then echo $x >&2; exit 99; fi
input:
apples grow in England
no output
error output:
apples grow in England
exit code: 99
`, "")

	st.Reset()
	c.CheckStderr(nil)
	c.Run(&st, "something")
	st.Expect(t, true, true, `unexpected error output
command: /bin/sh -c read x; if [ "$x" != nothing ]; then echo $x >&2; exit 99; fi
input:
something
no output
error output:
something
exit code: 99
`, "")

	st.Reset()
	c.Run(&st, "nothing")
	st.Expect(t, false, false, "", "")
}

func TestCmdCode(t *testing.T) {
	var st StubReporter
	c := Command("/bin/sh", "-c", "read x; exit $x")
	c.WantCode(17)

	c.Run(&st, "17")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "31")
	st.Expect(t, true, true, `incorrect exit code
command: /bin/sh -c read x; exit $x
input:
31
no output
no error output
exit code: 31
`, "")

	st.Reset()
	c.Run(&st, "0")
	st.Expect(t, true, true, `incorrect exit code
command: /bin/sh -c read x; exit $x
input:
0
no output
no error output
exit code: 0
`, "")

	st.Reset()
	c.CheckCode(func(actual int) bool {
		if big.NewInt(int64(actual)).ProbablyPrime(0) {
			return true
		}
		// We can report more information if we have access to the Reporter
		st.Helper()
		st.Errorf("%d is not prime", actual)
		return false
	})
	c.Run(&st, "17")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "31")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "99")
	st.Expect(t, true, true, `99 is not prime
incorrect exit code
command: /bin/sh -c read x; exit $x
input:
99
no output
no error output
exit code: 99
`, "")

	st.Reset()
	c.CheckCode(nil)
	c.Run(&st, "0")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c.Run(&st, "1")
	st.Expect(t, true, true, `non-zero exit code
command: /bin/sh -c read x; exit $x
input:
1
no output
no error output
exit code: 1
`, "")

	st.Reset()
	c2 := Command("/bin/sh", "-c", "echo oops >&2; read x; exit $x")
	c2.WantStderr("oops\n")
	c2.Run(&st, "3")
	st.Expect(t, false, false, "", "")

	st.Reset()
	c2.Run(&st, "0")
	st.Expect(t, true, true, `error output produced but exit code was 0
command: /bin/sh -c echo oops >&2; read x; exit $x
input:
0
no output
error output:
oops
exit code: 0
`, "")

	st.Reset()
	c2.WantStderr("hunky dory\n")
	c2.Run(&st, "0")
	st.Expect(t, true, true, `incorrect error output
command: /bin/sh -c echo oops >&2; read x; exit $x
input:
0
no output
error output:
oops
exit code: 0
`, "")

	st.Reset()
	c2.WantStdout("erewhon\n")
	c2.Run(&st, "0")
	st.Expect(t, true, true, `incorrect output
incorrect error output
command: /bin/sh -c echo oops >&2; read x; exit $x
input:
0
no output
error output:
oops
exit code: 0
`, "")

	st.Reset()
	c2.WantStderr("oops\n")
	c2.Run(&st, "0")
	st.Expect(t, true, true, `incorrect output
command: /bin/sh -c echo oops >&2; read x; exit $x
input:
0
no output
error output:
oops
exit code: 0
`, "")
}

func TestCmdPanic(t *testing.T) {
	var c Cmd
	msg := MustPanic(t, func() {
		c.Run(t, "")
	})
	Expect(t, "gotest.Cmd not initialized; use gotest.Command to create Cmds", msg.(string))
}

func TestCmdGo(t *testing.T) {
	c := Command("go", "version")
	c.WantStdout(fmt.Sprintf("go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	c.Run(t, "")
}

func TestCmdChdir(t *testing.T) {
	tmp := t.TempDir()

	path := filepath.Join(tmp, "somefile")
	f, e := os.Create(path)
	NilError(t, e)
	NilError(t, f.Close())

	c := Command("/bin/ls")
	c.Chdir(tmp)
	c.WantStdout("somefile\n")
	c.Run(t, "")

	var st StubReporter
	nondir := filepath.Join(tmp, "nondir")
	c.Chdir(nondir)
	c.Run(&st, "")
	if !st.Failed() {
		t.Error("running command in non-existent directory did not fail")
	}
	if !st.Killed() {
		t.Error("running command in non-existent directory did not stop test case")
	}
	if !strings.Contains(st.Logged(), nondir) {
		t.Error("bad error message for non-existent directory:", st.Logged())
	}
}
