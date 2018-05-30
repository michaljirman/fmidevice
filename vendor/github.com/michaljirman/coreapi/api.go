package coreapi

// API represent core informations
type API struct {
	BaseURL       string //https://httpbin.org
	DefaultRouter *CBRouter
	Client        *Client
}

// NewAPI create a new API struct
func NewAPI(baseURL string) *API {
	return &API{
		BaseURL:       baseURL,
		DefaultRouter: NewRouter(),
		Client:        NewClient(),
	}
}

// SetAuth sets auth on API
func (api *API) SetAuth(auth Authentication) {
	api.Client.SetAuth(auth)
}

func (api *API) SetDefaultHeaders(defaultHeaders map[string]string) {
	api.Client.SetDefaultHeaders(defaultHeaders)
}

// Call performs API call on specific resource and processes request through API's client
func (api *API) Call(resource *RestResource, params map[string]string,
	payload interface{}, additionalHeaders map[string]string) error {
	if err := api.Client.ProcessRequest(api.BaseURL, resource, params,
		payload, additionalHeaders); err != nil {
		return err
	}
	return nil
}
