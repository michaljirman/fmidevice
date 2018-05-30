package coreapi

import (
	"net/http"
	"testing"
)

func TestProcessRequest(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/plain")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/status/404/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/plain")
		w.WriteHeader(http.StatusNotFound)
	})

	client := NewClient()
	router := NewRouter()
	router.RegisterFunc(200, func(resp *http.Response) error {
		return nil
	})
	router.RegisterFunc(404, func(resp *http.Response) error {
		return nil
	})
	resource := NewResource("/get", "GET", router)
	if err := client.ProcessRequest(server.URL+"/status/404", resource, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := client.ProcessRequest(server.URL, resource, nil, nil, nil); err != nil {
		t.Fatal(err)
	}

}
