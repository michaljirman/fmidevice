package fmip

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
	accountName = "tester.test@icloud.com"
	password    = "tester2018"
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

func TestInitAndRefreshMobileClientAPICall(t *testing.T) {
	teardown := setup()
	defer teardown()

	resourceInitClientURL := fmt.Sprintf("/fmipservice/device/%s/initClient", accountName)
	mux.HandleFunc(resourceInitClientURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("init_mobile_client_response.json"))
	})

	resourceRefreshURL := fmt.Sprintf("/fmipservice/device/%s/refreshClient", accountName)
	mux.HandleFunc(resourceRefreshURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("refresh_mobile_client_response.json"))
	})

	fmipService, err := NewFmipMobileService(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	fmipService.InitMobileClient(accountName, password)
	prsID := fmipService.FmipClientResponse.ServerContext.PrsID
	authToken := fmipService.FmipClientResponse.ServerContext.AuthToken

	if fmipService.FmipClientResponse == nil {
		t.Fatal()
	}
	if prsID == 0 {
		t.Fatal()
	}
	if authToken == "" {
		t.Fatal()
	}

	fmipService.RefreshMobileClient(accountName, prsID, authToken)

	if fmipService.FmipClientResponse == nil {
		t.Fatal(err)
	}
	if prsID == 0 {
		t.Fatal()
	}
	if authToken == "" {
		t.Fatal()
	}

	if len(fmipService.FmipClientResponse.Content) != 1 {
		t.Fatal()
	}
}
