// Copyright 2023 Patrick Smith
// Use of this source code is subject to the MIT-style license in the LICENSE file.

package gotest

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRequire(t *testing.T) {
	var st StubReporter
	Require(&st, true)
	st.Expect(t, false, false, "")

	Require(&st, false)
	st.Expect(t, true, true, "Required condition failed\n")
}

func TestExpect(t *testing.T) {
	var st StubReporter
	Expect(&st, 5, 5)
	st.Expect(t, false, false, "")

	Expect(&st, "a", "b")
	st.Expect(t, true, true, "Expected a but actual value was b\n")

	// This should not compile, as the arguments have different types: Expect(&st, 7, "7")
	testprogram := `package foo
import "github.com/pat42smith/gotest"
import "testing"

func Try(t *testing.T) {
	gotest.Expect(t, 7, "7")
}`

	tmp := t.TempDir()
	testfile := filepath.Join(tmp, "foo.go")
	if e := os.WriteFile(testfile, []byte(testprogram), 0444); e != nil {
		t.Fatal(e)
	}

	cmd := Command("go", "build", testfile)
	cmd.CheckStderr(func(actual string) bool {
		return strings.Contains(actual, `mismatched types untyped int and untyped string (cannot infer T)`)
	})
	cmd.Run(t, "")
}

func TestPanics(t *testing.T) {
	p, w := panics(func() {})
	Expect(t, false, p)
	Require(t, w == nil) // Can't use Expect(t, nil, w) because w (type any) does not implement comparable.

	p, w = panics(func() {
		panic(97)
	})
	Expect(t, true, p)
	Require(t, w == 97)

	// panic(nil) is a special case that, in older versions of Go, appeared to the recover
	// function as if no panic happened.  Since Go 1.21, it causes a runtime panic.
	p, w = panics(func() {
		panic(nil)
	})
	Expect(t, true, p)
	if _, ok := w.(*runtime.PanicNilError); !ok {
		t.Fatalf("result of panics after panic(nil) is a %T; expected *runtime.PanicNilError", w)
	}
}

func TestMustPanic(t *testing.T) {
	var st StubReporter
	x := MustPanic(&st, func() {
		panic("oops")
	})
	st.Expect(t, false, false, "")
	Require(t, x == "oops")

	x = MustPanic(&st, func() {})
	st.Expect(t, true, true, "Expected panic did not occur\n")
	Require(t, x == nil)
}

func TestNotFatal(t *testing.T) {
	var st1, st2, st3 StubReporter
	NotFatal{&st1}.FailNow()
	st1.Expect(t, true, false, "")

	NotFatal{&st2}.Fatal("problem")
	st2.Expect(t, true, false, "problem\n")

	NotFatal{&st3}.Fatalf("<%s>", "uh oh")
	st3.Expect(t, true, false, "<uh oh>\n")
}
