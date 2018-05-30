package coreapi

import (
	"net/http"
	"net/url"
	"testing"
)

func TestUnknownStatusCode(t *testing.T) {
	router := NewRouter()
	fakeURL, err := url.Parse("https://httpbin.org/doesnotexist")
	if err != nil {
		t.Fatal(err)
	}
	resp := &http.Response{
		Request: &http.Request{
			URL: fakeURL,
		},
		StatusCode: 404,
	}

	if err := router.CallFunc(resp); err == nil {
		t.Fatal(err)
	}
}
