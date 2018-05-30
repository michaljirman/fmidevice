package idmsa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

const (
	clientID               = "580453d1-1fa3-4d20-bf80-a9faa828ef53"
	clientBuildNumber      = "1809Project50"
	clientMasteringNumber  = "1809B29"
	xAppleSessionToken     = "0WvxHOC4Aw0N25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUP+LtgN25ljk45jkl23jlk+235jkl5+jkl2jl254k;25j4klj5j5l432j5324u0cAABMNvUPfdasdsfa+Ltg="
	xAppleIDSessionID      = "5F240D72D67FDA679SA967F679A93517"
	xAppleIDAccountCountry = "GBR"
	accountLoginUIResource = "/appleauth/auth/signin"
	accountName            = "tester.test@icloud.com"
	password               = "tester2018"
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	return func() {
		server.Close()
	}
}
func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestIdmsaSigninAPICall(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/appleauth/auth/signin", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("X-Apple-ID-Session-Id", xAppleIDSessionID)
		w.Header().Set("X-Apple-Session-Token", xAppleSessionToken)
		w.Header().Set("X-Apple-ID-Account-Country", xAppleIDAccountCountry)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, fixture("signin_response.json"))
	})

	idmsaService, err := NewIdmsaService(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	idmsaService.AccountLoginUI(accountLoginUIResource, accountName, password)

	if idmsaService.XappleSessionToken == "" || idmsaService.XappleSessionToken != xAppleSessionToken {
		t.Fatal()
	}

	if idmsaService.XappleIDSessionID == "" || idmsaService.XappleIDSessionID != xAppleIDSessionID {
		t.Fatal()
	}

	if idmsaService.XappleIDAccountCountry == "" || idmsaService.XappleIDAccountCountry != xAppleIDAccountCountry {
		t.Fatal()
	}
}
