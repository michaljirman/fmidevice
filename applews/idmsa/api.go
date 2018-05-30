package idmsa

import "github.com/michaljirman/coreapi"

var defaultHeaders = map[string]string{
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36",
	"Content-Type": "application/json",
	"Accept":       "application/json, text/javascript, */*; q=0.01",
	"Connection":   "keep-alive",
	//"Accept-Encoding": "gzip, deflate, br",
	"Accept-Language":  "en-US,en;q=0.9,cs;q=0.8",
	"Origin":           "https://idmsa.apple.com",
	"X-Requested-With": "XMLHttpRequest",
}

func NewIdmsaAPI(baseURL string) *coreapi.API {
	idmsaAPI := coreapi.NewAPI(baseURL)
	idmsaAPI.SetDefaultHeaders(defaultHeaders)
	return idmsaAPI
}
