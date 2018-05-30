package coreapi

import (
	"fmt"
	"net/http"
)

// RouterFunc is the callback function type
type RouterFunc func(resp *http.Response) error

// CBRouter represents a collection of routers based on status codes
type CBRouter struct {
	Routers       map[int]RouterFunc
	DefaultRouter RouterFunc
}

// NewRouter returns a new router
func NewRouter() *CBRouter {
	return &CBRouter{
		Routers: make(map[int]RouterFunc),
		DefaultRouter: func(resp *http.Response) error {
			return fmt.Errorf("From: %s received unknown status: %d",
				resp.Request.URL.String(), resp.StatusCode)
		},
	}
}

// RegisterFunc will register a function with a status code
func (r *CBRouter) RegisterFunc(status int, fn RouterFunc) {
	r.Routers[status] = fn
}

// CallFunc calls a registered function in the router
func (r *CBRouter) CallFunc(resp *http.Response) error {
	fn, ok := r.Routers[resp.StatusCode]
	if !ok {
		fn = r.DefaultRouter
	}
	if err := fn(resp); err != nil {
		return err
	}
	return nil
}

// DefaultRouter returns an error containing http status code
func DefaultRouter(resp *http.Response) error {
	return fmt.Errorf("status code %d", resp.StatusCode)
}
