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
