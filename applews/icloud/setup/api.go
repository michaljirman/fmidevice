package setup

import "github.com/michaljirman/coreapi"

var defaultHeaders = map[string]string{
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36",
	"Content-Type": "text/plain",
	"Accept":       "application/json, text/javascript, */*; q=0.01",
	"Connection":   "keep-alive",
	//"Accept-Encoding": "gzip, deflate, br",
	"Accept-Language": "en-US,en;q=0.9,cs;q=0.8",
	"Origin":          "https://www.icloud.com",
}

func NewSetupAPI(baseURL string) *coreapi.API {
	setupAPI := coreapi.NewAPI(baseURL)
	setupAPI.SetDefaultHeaders(defaultHeaders)
	return setupAPI
}
