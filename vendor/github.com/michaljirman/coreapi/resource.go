package coreapi

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
)

// RestResource for API
type RestResource struct {
	Endpoint string // /get
	Method   string // GET
	Router   *CBRouter
}

// NewResource returns a new RestResource
func NewResource(endpoint, method string, router *CBRouter) *RestResource {
	return &RestResource{
		Endpoint: endpoint,
		Method:   method,
		Router:   router,
	}
}

// RenderEndpoint processes an url and handles its params.
func (r *RestResource) RenderEndpoint(params map[string]string) string {
	if params == nil {
		return r.Endpoint
	}
	t, err := template.New("resource").Parse(r.Endpoint)
	if err != nil {
		log.Fatalln("Unable to parse endpoint")
	}
	buffer := &bytes.Buffer{}
	t.Execute(buffer, params)
	endpoint, err := ioutil.ReadAll(buffer)
	if err != nil {
		log.Fatalln("Unable to read endpoint")
	}
	return string(endpoint)
}
