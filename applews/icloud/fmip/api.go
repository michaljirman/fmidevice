package fmip

import "github.com/michaljirman/coreapi"

var fmipMobileDefaultHeaders = map[string]string{
	"Content-Type": "text/plain",
	"Accept":       "application/json, text/javascript, */*; q=0.01",
	"Connection":   "keep-alive",
	//"Accept-Encoding": "gzip, deflate, br",
	"Accept-Language":       "en-US,en;q=0.9,cs;q=0.8",
	"Origin":                "https://www.icloud.com",
	"X-Apple-Realm-Support": "1.0",
	"X-Apple-Find-API-Ver":  "3.0",
	"User-Agent":            "FindMyiPhone/500 CFNetwork/758.4.3 Darwin/15.5.0",
}

func NewFmipMobileAPI(baseURL string) *coreapi.API {
	api := coreapi.NewAPI(baseURL)
	api.SetDefaultHeaders(fmipMobileDefaultHeaders)
	return api
}

var fmipWebDefaultHeaders = map[string]string{
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36",
	"Content-Type": "text/plain",
	"Accept":       "application/json, text/javascript, */*; q=0.01",
	"Connection":   "keep-alive",
	//"Accept-Encoding": "gzip, deflate, br",
	"Accept-Language": "en-US,en;q=0.9,cs;q=0.8",
	"Origin":          "https://www.icloud.com",
}

func NewFmipWebAPI(baseURL string) *coreapi.API {
	// https://p66-fmipweb.icloud.com
	fmipWebAPI := coreapi.NewAPI(baseURL)
	fmipWebAPI.SetDefaultHeaders(fmipWebDefaultHeaders)
	return fmipWebAPI
}
