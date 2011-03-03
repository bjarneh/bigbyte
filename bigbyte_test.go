// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// I can rip off tests from the standard library I guess
// as long as i keep that header, I think that's what
// BSD means, but I'm not sure... :-)

package bigbyte_test

import (
    . "github.com/bjarneh/bigbyte"
    "testing"
    "io/ioutil"
    "bytes"
    "os"
    "path"
)

var haystack    []byte

func init(){

    var err os.Error
    var srcroot, me string

    srcroot = os.Getenv("SRCROOT")
    me  = "github.com/bjarneh/bigbyte"

    // we are testing with godag
    if srcroot != "" {
        haystack, err = ioutil.ReadFile(path.Join(srcroot,me,"wild-duck-no-utf8.txt"))
    }else{// we are testing with Makefile
        haystack, err = ioutil.ReadFile("wild-duck-no-utf8.txt")
    }

    if err != nil {
        panic("could not read file needed for Benchmark")
    }
}

type BinOpTest struct {
	a string
	b string
	i int
}

var indexTests = []BinOpTest{
	{"", "", 0},
	{"", "a", -1},
	{"", "foo", -1},
	{"fo", "foo", -1},
	{"foo", "foo", 0},
	{"oofofoofooo", "f", 2},
	{"oofofoofooo", "foo", 4},
	{"barfoobarfoo", "foo", 3},
	{"foo", "", 0},
	{"foo", "o", 1},
	{"abcABCabc", "A", 3},
	// cases with one byte strings - test IndexByte and special case in Index()
	{"", "a", -1},
	{"x", "a", -1},
	{"x", "x", 0},
	{"abc", "a", 0},
	{"abc", "b", 1},
	{"abc", "c", 2},
	{"abc", "x", -1},
	{"barfoobarfooyyyzzzyyyzzzyyyzzzyyyxxxzzzyyy", "x", 33},
	{"foofyfoobarfoobar", "y", 4},
	{"oooooooooooooooooooooo", "r", -1},
}


// Execute f on each test case.  funcName should be the name of f; it's used
// in failure reports.
func runIndexTests(t *testing.T, f func(s, sep []byte) int, funcName string, testCases []BinOpTest) {
	for _, test := range testCases {
		a := []byte(test.a)
		b := []byte(test.b)
		actual := f(a, b)
		if actual != test.i {
			t.Errorf("%s(%q,%q) = %v; want %v", funcName, a, b, actual, test.i)
		}
	}
}

func TestIndexBMH(t *testing.T)  { runIndexTests(t, IndexBMH, "IndexBMH", indexTests) }

// using the same philosophy as the one the Go developers used for testing,
// only now for benching :-)

type BenchTest struct {
    find string
    where int
}

var benchTests = []BenchTest{
	{"lkjlxkcj lsdjfl s", -1},
	{"strækker armene ud og mumler.  Lovet være herren;", 276159},
	{"her viser jo hjertelag,", 16573},
	{"I alt det strå?  går mod den øverste dør til", 58757},
	{"LKJoisdjflksjdf iOJJDFjl skdfjls df", -1},
}


func BenchmarkIndexBMH(b *testing.B){
    runBenchTests(b, IndexBMH, "bigbyte.IndexBMH", benchTests)
}

func BenchmarkIndex(b *testing.B){
    runBenchTests(b, bytes.Index, "bytes.Index", benchTests)
}

func runBenchTests(b *testing.B, f func(s, sep []byte) int, funcName string, testCases []BenchTest){

    var needle []byte
    var result int

    for _, bt := range testCases {
        needle = []byte(bt.find)
        result = f(haystack, needle)
        if result != bt.where {
            panic("function: "+funcName+", does not work properly")
        }
        for i := 0; i < b.N; i++ {
            _ = f(haystack, needle)
        }
    }
}

