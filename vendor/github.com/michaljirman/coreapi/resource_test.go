package coreapi

import "testing"

func TestEndpointTemplate(t *testing.T) {
	resource := &RestResource{
		Endpoint: "/user/{{.user}}",
		Method:   "GET",
		Router:   NewRouter(),
	}
	renderedEndpoint := resource.RenderEndpoint(map[string]string{
		"user": "tester",
	})
	if renderedEndpoint != "/user/tester" {
		t.Fail()
	}
}
