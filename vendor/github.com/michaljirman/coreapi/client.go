package coreapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Client represents an interlan api client wrapping standard http.Client and Authentication
type Client struct {
	Client         *http.Client
	AuthInfo       Authentication
	DefaultHeaders map[string]string
}

// NewClient retursn a new client
func NewClient() *Client {
	return &Client{
		// Client: http.DefaultClient,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// SetAuth will set AuthInfo
func (c *Client) SetAuth(auth Authentication) {
	c.AuthInfo = auth
}

// SetAuth will set AuthInfo
func (c *Client) SetDefaultHeaders(defaultHeaders map[string]string) {
	c.DefaultHeaders = defaultHeaders
}

// ProcessRequest performs the api request and handles response
func (c *Client) ProcessRequest(baseURL string, resource *RestResource, params map[string]string,
	payload interface{}, additionalHeaders map[string]string) error {
	endpoint := strings.TrimLeft(resource.RenderEndpoint(params), "/")
	trimmedBAseURL := strings.TrimRight(baseURL, "/")
	url := trimmedBAseURL + "/" + endpoint
	req := buildClientRequst(url, resource.Method, payload)
	//Set default headers
	for k, v := range c.DefaultHeaders {
		req.Header.Set(k, v)
	}
	//Set additional headers
	for k, v := range additionalHeaders {
		req.Header.Set(k, v)
	}

	if c.AuthInfo != nil {
		req.Header.Add("Authorization", c.AuthInfo.AuthorizationHeader())
	}

	resp, err := c.Client.Do(req)
	// fmt.Println(resp)
	if err != nil {
		return err
	}
	return resource.Router.CallFunc(resp)
}

func buildClientRequst(url, method string, payload interface{}) *http.Request {
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil
		}
		payloadBuffer := bytes.NewBuffer(payloadBytes)
		req, err := http.NewRequest(method, url, payloadBuffer)
		return req
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil
	}
	return req
}
