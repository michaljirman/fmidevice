package coreapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	api    *API
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	api = NewAPI(server.URL)

	return func() {
		server.Close()
	}
}

func TestAPICall(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/plain")
		w.WriteHeader(http.StatusOK)
	})

	router := NewRouter()
	router.RegisterFunc(200, func(resp *http.Response) error {
		return nil
	})
	resource := NewResource("/get", "GET", router)
	if err := api.Call(resource, nil, nil, nil); err != nil {
		t.Fatal(err)
	}

}
