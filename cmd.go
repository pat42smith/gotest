// Copyright 2023 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"os/exec"
	"strings"
)

// A Cmd runs an external command inside a test case
// and checks the results of the command.
type Cmd struct {
	name               string
	args               []string
	dir                string
	checkOut, checkErr func(actual string) bool
	checkCode          func(actual int) bool
}

// Command creates a Cmd object to run a specific command once.
//
// The arguments name and args are the same as for os/exec.Command.
//
// As of go 1.19, go test puts GOROOT/bin at the beginning of the PATH
// when running test cases, so it is appropriate to use "go" for name
// in order to invoke the go command inside a test case. However,
// if you build the test case executable and run it separately,
// you will be responsible for ensuring the go command is in your PATH.
func Command(name string, args ...string) *Cmd {
	var cmd Cmd
	cmd.name = name
	cmd.args = args
	return &cmd
}

// CheckStdout sets the function used to check the command's output.
//
// The check function will be passed the output produced by the command,
// and should return whether the output is correct.
//
// CheckStdout(nil), the default, is equivalent to
// CheckStdout(func (actual string) bool { return actual == "" }).
func (c *Cmd) CheckStdout(check func(actual string) bool) {
	c.checkOut = check
}

// CheckStderr sets the function used to check the error output produced by the command.
//
// The check function will be passed the error output produced by the command,
// and should return whether the error output is correct.
//
// CheckStderr(nil), the default, is equivalent to
// CheckStderr(func (actual string) bool { return actual == "" }).
func (c *Cmd) CheckStderr(check func(actual string) bool) {
	c.checkErr = check
}

// CheckCode sets the function used to check the command's exit code.
//
// The check function will be passed the command's exit code, and
// should return whether the code is correct.
//
// CheckCode(nil), the default, has two effects. If either the stdout
// or the stderr from the command is deemed incorrect, then the exit
// code is not checked at all. If the code is checked, then it is
// expected to be 0 if the command produced no error output, and non-0
// otherwise.
func (c *Cmd) CheckCode(check func(actual int) bool) {
	c.checkCode = check
}

// WantStdout indicates that the output of the command should be exactly expected.
func (c *Cmd) WantStdout(expected string) {
	c.checkOut = func(actual string) bool {
		return actual == expected
	}
}

// WantStderr indicates that the error output of the command should be exactly expected.
func (c *Cmd) WantStderr(expected string) {
	c.checkErr = func(actual string) bool {
		return actual == expected
	}
}

// WantCode indicates that the exit code of the command should be expected.
func (c *Cmd) WantCode(expected int) {
	c.checkCode = func(actual int) bool {
		return actual == expected
	}
}

// Chdir sets the working directory where the command will be run.
// Chdir(""), the default, is equivalent to Chdir("."); it uses
// the current directory.
func (c *Cmd) Chdir(path string) {
	c.dir = path
}

// Run runs the external command and checks the results.
//
// The content of input is passed to the command as its stdin.
// The results of the command are checked per previous calls to
// the Check* and Want* methods; any test failures are reported to t.
//
// If there are any failures, Run will record through t the command
// executed, its output, error output, and exit code. It will then
// call t.FailNow.
//
// If the command can not be started or is terminated by a signal,
// Run will report a fatal error and skip checking the command results.
//
// It is permissible to call Run multiple times on the same Cmd object,
// in order to test the same external command with varying inputs.
// The Check* or Want* functions may be called between calls to Run,
// if the expected results will change.
func (c *Cmd) Run(t Reporter, input string) {
	t.Helper()
	if c.name == "" {
		panic("gotest.Cmd not initialized; use gotest.Command to create Cmds")
	}

	cmd := exec.Command(c.name, c.args...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Dir = c.dir

	var out, err strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &err
	e := cmd.Run()

	code := 0
	if e != nil {
		ee, ok := e.(*exec.ExitError)
		if ok {
			code = ee.ExitCode()
			ok = ee.Exited()
		}
		if !ok {
			t.Fatal(e)
			return // In case t.Fatal has been overridden to not terminate the test case.
		}
	}

	ok := true

	if c.checkOut == nil {
		if out.Len() > 0 {
			t.Error("unexpected output")
			ok = false
		}
	} else if !c.checkOut(out.String()) {
		t.Error("incorrect output")
		ok = false
	}

	if c.checkErr == nil {
		if err.Len() > 0 {
			t.Error("unexpected error output")
			ok = false
		}
	} else if !c.checkErr(err.String()) {
		t.Error("incorrect error output")
		ok = false
	}

	if c.checkCode == nil {
		if ok {
			if err.Len() == 0 {
				if code != 0 {
					t.Error("non-zero exit code")
					ok = false
				}
			} else {
				if code == 0 {
					t.Error("error output produced but exit code was 0")
					ok = false
				}
			}
		}
	} else if !c.checkCode(code) {
		t.Error("incorrect exit code")
		ok = false
	}

	if !ok {
		if len(c.args) == 0 {
			t.Errorf("command: %s", c.name)
		} else {
			t.Errorf("command: %s %s", c.name, strings.Join(c.args, " "))
		}
		if len(input) == 0 {
			t.Error("no input")
		} else {
			// Not t.Error(...), in case the input ends with a newline.
			t.Errorf("input:\n%s", input)
		}
		if out.Len() == 0 {
			t.Error("no output")
		} else {
			// Don't use t.Error("output:\n" + out.String()); the output usually ends with a newline,
			// and t.Error always adds another newline.
			t.Errorf("output:\n%s", out.String())
		}
		if err.Len() == 0 {
			t.Error("no error output")
		} else {
			// Again not using t.Error
			t.Errorf("error output:\n%s", err.String())
		}
		t.Errorf("exit code: %d", code)
		t.FailNow()
	}
}
