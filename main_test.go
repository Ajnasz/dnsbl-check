package main

import (
	"testing"
)

func Test_negate(t *testing.T) {
	if negate(isEmptyString)("") {
		t.Errorf("It should be false")
	}
	if !negate(isEmptyString)("a") {
		t.Errorf("It should be true")
	}
}

func Test_getToken(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "",
			Expected: "",
		},
		{
			Input:    "foo",
			Expected: "foo",
		},
		{
			Input:    "foo bar",
			Expected: "foo",
		},
	}

	for _, testCase := range testCases {
		actual := getToken(testCase.Input, " ")

		if actual != testCase.Expected {
			t.Errorf("expected %q got %q", testCase.Expected, actual)
		}
	}
}

func Benchmark_getToken(b *testing.B) {
	noop := func(string) {}
	for i := 0; i < b.N; i++ {
		str := getToken("foo bar", " ")
		noop(str)
	}
}
