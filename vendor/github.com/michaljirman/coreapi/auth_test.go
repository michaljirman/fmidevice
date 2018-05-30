package coreapi

import "testing"

func TestAuthBasic(t *testing.T) {
	basic := NewAuthBasic("user", "password")
	header := basic.AuthorizationHeader()
	if header != "Basic dXNlcjpwYXNzd29yZA==" {
		t.Fail()
	}
}
