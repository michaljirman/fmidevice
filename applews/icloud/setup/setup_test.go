package setup

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
	accountLoginResource   = "/setup/ws/1/accountLogin"

	xAppleWebAuthValidateCookie = `X-APPLE-WEBAUTH-LOGIN="v=1:t=IAAAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABLwIAAAAAFhlCqMQ~~";Path=/;Domain=.icloud.com;Secure;HttpOnly`
	xAppleWebAuthLoginCookie    = `X-APPLE-WEBAUTH-VALIDATE="v=1:t=IAAAAAAABLwIAAAAAFsOYAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABBUpEVk5lA~~";Path=/;Domain=.icloud.com;Secure`
	cookieResult                = `X-APPLE-WEBAUTH-VALIDATE="v=1:t=IAAAAAAABLwIAAAAAFsOYAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABAAAAABBUpEVk5lA~~";Path=/;Domain=.icloud.com;Secure;`
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

func TestSetupValidateAPICall(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/setup/ws/1/validate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/plain")
		w.WriteHeader(421)
		fmt.Fprint(w, fixture("validate_response.json"))
	})

	ss, err := NewSetupService(server.URL, clientBuildNumber, clientID, clientMasteringNumber)

	if err != nil {
		t.Fatal(err)
	}
	ss.Validate()
	if ss.SetupRequestResponse == nil {
		t.Fatal()
	}
	if ss.SetupRequestResponse.ConfigBag.Urls.AccountLogin == "" {
		t.Fatal()
	}
	// t.Log(ss.SetupRequestResponse.ConfigBag.Urls.AccountLoginUI)
	if ss.SetupRequestResponse.ConfigBag.Urls.AccountLoginUI == "" {
		t.Fatal()
	}
}

func TestSetupLoginAPICall(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc(accountLoginResource,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			w.Header().Set("Set-Cookie", xAppleWebAuthValidateCookie)
			w.Header().Set("Set-Cookie", xAppleWebAuthLoginCookie)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, fixture("account_login_response.json"))
		})

	ss, err := NewSetupService(server.URL, clientBuildNumber, clientID, clientMasteringNumber)

	if err != nil {
		t.Fatal(err)
	}
	ss.AccountLogin(accountLoginResource, xAppleSessionToken, xAppleIDAccountCountry)

	if ss.AccountLoginResponse == nil {
		t.Fatal()
	}

	if ss.AccountLoginResponse.DsInfo.Dsid == "" {
		t.Fatal()
	}

	if ss.AccountLoginResponse.Webservices.Findme.URL == "" {
		t.Fatal()
	}

	xAppleCookiesHeader := ss.GetXappleCookiesHeader()
	t.Log(xAppleCookiesHeader)
	if xAppleCookiesHeader != cookieResult {
		t.Fatal()
	}

}
