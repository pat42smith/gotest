// Copyright 2022 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRequire(t *testing.T) {
	var st StubReporter
	Require(&st, true)
	st.Expect(t, false, false, "", "after Require(true)")

	Require(&st, false)
	st.Expect(t, true, true, "Test requirement failed\n", "after Require(false)")
}

func TestNilError(t *testing.T) {
	var st StubReporter
	NilError(&st, nil)
	st.Expect(t, false, false, "", "after NilError(nil)")

	NilError(&st, fmt.Errorf("oops"))
	st.Expect(t, true, true, "oops\n", "after NilError(oops)")
}

func TestExpect(t *testing.T) {
	var st StubReporter
	Expect(&st, 5, 5)
	st.Expect(t, false, false, "", "after Expect(5,5)")

	Expect(&st, "a", "b")
	st.Expect(t, true, true, "Expected a but actual value was b\n", "after Expect(a,b)")

	// This should not compile, as the arguments have different types: Expect(&st, 7, "7")
	testprogram := `package main_test
import "github.com/pat42smith/gotest"
import "testing"

func TestExpect(t *testing.T) {
	gotest.Expect(t, 7, "7")
}`

	here, e := os.Getwd()
	NilError(t, e)
	here = strconv.Quote(here)
	gomod := `module main
go 1.18
require github.com/pat42smith/gotest v0.0.0
replace github.com/pat42smith/gotest => ` + here + "\n"

	tmp := t.TempDir()
	testfile := filepath.Join(tmp, "main_test.go")
	NilError(t, os.WriteFile(testfile, []byte(testprogram), 0444))
	modfile := filepath.Join(tmp, "go.mod")
	NilError(t, os.WriteFile(modfile, []byte(gomod), 0444))

	gocmd, e := exec.LookPath("go")
	NilError(t, e)
	cmd := exec.Command(gocmd, "test")
	cmd.Dir = tmp
	out, e := cmd.CombinedOutput()
	if e == nil {
		t.Fatal("Expected error from go test did not happen")
	}
	need := []byte(`default type string of "7" does not match inferred type int for T`)
	Require(t, bytes.Contains(out, need))
}

func Testpanics(t *testing.T) {
	p, w := panics(func() {})
	Expect(t, false, p)
	Require(t, w == nil) // Can't use Expect(t, nil, w) because w (type any) does not implement comparable.

	p, w = panics(func() {
		panic(97)
	})
	Expect(t, true, p)
	Require(t, w == 97)

	// panic(nil) is a special case that appears to the recover function as if no panic happened.
	// So be sure to check this case.
	p, w = panics(func() {
		panic(nil)
	})
	Expect(t, true, p)
	Require(t, w == nil)
}

func TestMustPanic(t *testing.T) {
	var st StubReporter
	x := MustPanic(&st, func() {
		panic("oops")
	})
	st.Expect(t, false, false, "", "MustPanic(panic)")
	Require(t, x == "oops")

	x = MustPanic(&st, func() {})
	st.Expect(t, true, true, "Expected panic did not occur\n", "MustPanic()")
	Require(t, x == nil)
}

func TestErrorsNotFatal(t *testing.T) {
	var st1, st2, st3 StubReporter
	ErrorsNotFatal{&st1}.FailNow()
	st1.Expect(t, true, false, "", "after ErrorsNotFatal.FailNow")

	ErrorsNotFatal{&st2}.Fatal("problem")
	st2.Expect(t, true, false, "problem\n", "after ErrorsNotFatal.Fatal")

	ErrorsNotFatal{&st3}.Fatalf("<%s>", "uh oh")
	st3.Expect(t, true, false, "<uh oh>\n", "after ErrorsNotFatal.Fatalf")
}
