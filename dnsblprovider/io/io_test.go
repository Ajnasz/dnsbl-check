package io

import "testing"

type provider struct {
	retval bool
	err    error
	name   string
	reason string
}

func (p provider) IsBlacklisted(adress string) (bool, error) {
	return p.retval, p.err
}

func (p provider) GetName() string {
	return p.name
}

func (p provider) GetReason(string) (string, error) {
	return p.reason, nil
}

func Test_Lookup(t *testing.T) {
	p := provider{
		retval: true,
	}
	ret := Lookup("127.0.0.1", p)

	if ret.IsBlacklisted == false {
		t.Errorf("Expected to be false")
	}
}
