package providerlist

import (
	"testing"
)

func Test_negate(t *testing.T) {
	if negate(isEmptyString)("") != false {
		t.Errorf("It should be false")
	}

	if negate(isEmptyString)("a") != true {
		t.Errorf("It should be true")
	}
}
