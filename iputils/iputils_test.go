package iputils

import "testing"

func Test_ReverseStringByToken(t *testing.T) {
	testCases := []struct {
		Input     string
		Expected  string
		Delimiter string
	}{
		{
			Input:     "",
			Expected:  "",
			Delimiter: " ",
		},
		{
			Input:     "foo",
			Expected:  "foo",
			Delimiter: " ",
		},
		{
			Input:     "foo bar",
			Expected:  "bar foo",
			Delimiter: " ",
		},
		{
			Input:     "foo bar baz",
			Expected:  "baz bar foo",
			Delimiter: " ",
		},
		{
			Input:     "1.2.3.4",
			Expected:  "4.3.2.1",
			Delimiter: ".",
		},
	}

	for _, testCase := range testCases {
		actual := reverseStringByToken(testCase.Input, testCase.Delimiter)

		if actual != testCase.Expected {
			t.Errorf("expected %q got %q", testCase.Expected, actual)
		}
	}
}

func Benchmark_ReverseStringByToken(b *testing.B) {
	noop := func(string) {}
	for i := 0; i < b.N; i++ {
		str := reverseStringByToken("foo bar", " ")
		noop(str)
	}
}
