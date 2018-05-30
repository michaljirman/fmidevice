package fmip

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/michaljirman/coreapi"
)

type FmipService struct {
	FmipClientResponse *fmipClientResponse
	api                *coreapi.API
}

func NewFmipMobileService(baseURL string) (*FmipService, error) {
	return &FmipService{
		FmipClientResponse: nil,
		api:                NewFmipMobileAPI(baseURL),
	}, nil
}

func NewFmipWebService(baseURL string) (*FmipService, error) {
	return &FmipService{
		FmipClientResponse: nil,
		api:                NewFmipWebAPI(baseURL),
	}, nil
}

type fmipClientResponse struct {
	ServerContext struct {
		AuthToken string `json:"authToken"`
		PrsID     int    `json:"prsId"`
	}
	Content []struct {
		Name              string  `json:"name"`
		DeviceDisplayName string  `json:"deviceDisplayName"`
		BatteryLevel      float32 `json:"batteryLevel"`
		BatteryStatus     string  `json:"batteryStatus"`
		Location          struct {
			Longitude float32 `json:"longitude"`
			Latitude  float32 `json:"latitude"`
		} `json:"location"`
	} `json:"content"`
}

func (fmipService *FmipService) clientSuccess(resp *http.Response) error {
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(content))
	fmipService.FmipClientResponse = &fmipClientResponse{}
	json.Unmarshal(content, fmipService.FmipClientResponse)
	return nil
}

func (fmipService *FmipService) InitMobileClient(accountName, password string) {
	router := coreapi.NewRouter()
	router.DefaultRouter = coreapi.DefaultRouter
	router.RegisterFunc(200, fmipService.clientSuccess)
	fmipService.api.SetAuth(coreapi.NewAuthBasic(accountName, password))
	resource := coreapi.NewResource("/fmipservice/device/{{.accountName}}/initClient", "POST", router)
	additionalHeaders := map[string]string{
		"X-Apple-AuthScheme": "UserIDGuest",
	}
	err := fmipService.api.Call(resource, map[string]string{
		"accountName": accountName,
	}, nil, additionalHeaders)
	if err != nil {
		log.Fatalln(err)
	}
}

func (fmipService *FmipService) RefreshMobileClient(accountName string, prsId int, authToken string) {
	router := coreapi.NewRouter()
	router.DefaultRouter = coreapi.DefaultRouter
	router.RegisterFunc(200, fmipService.clientSuccess)
	fmipService.api.SetAuth(coreapi.NewAuthBasic(strconv.Itoa(prsId), authToken))
	resource := coreapi.NewResource("/fmipservice/device/{{.accountName}}/refreshClient", "POST", router)
	additionalHeaders := map[string]string{
		"X-Apple-AuthScheme": "Forever",
	}
	err := fmipService.api.Call(resource, map[string]string{
		"accountName": accountName,
	}, nil, additionalHeaders)
	if err != nil {
		log.Fatalln(err)
	}
}

// {"clientContext":{"appName":"iCloud Find (Web)","appVersion":"2.0","timezone":"US/Pacific","inactiveTime":2026,"apiVersion":"3.0","deviceListVersion":1,"fmly":true}}
type clientContext struct {
	AppName           string `json:"appName"`
	AppVersion        string `json:"appVersion"`
	Timezone          string `json:"timezone"`
	InactiveTime      int    `json:"inactiveTime"`
	APIVersion        string `json:"apiVersion"`
	DeviceListVersion int    `json:"deviceListVersion"`
	Fmly              bool   `json:"fmly"`
}

type initClientRequestPayload struct {
	ClientContext clientContext `json:"clientContext"`
}

func (fmipService *FmipService) InitWebClient(xAppleCookiesHeader,
	clientBuildNumber, clientID, clientMasteringNumber, dsid string) {
	router := coreapi.NewRouter()
	router.DefaultRouter = coreapi.DefaultRouter
	router.RegisterFunc(200, fmipService.clientSuccess)
	resource := coreapi.NewResource(
		"/fmipservice/client/web/initClient?clientBuildNumber={{.clientBuildNumber}}&clientId={{.clientID}}&clientMasteringNumber={{.clientMasteringNumber}}&dsid={{.dsid}}",
		"POST", router)

	additionalHeaders := map[string]string{
		"Cookie": xAppleCookiesHeader,
	}

	payload := &initClientRequestPayload{
		ClientContext: clientContext{
			AppName:           "iCloud Find (Web)",
			AppVersion:        "2.0",
			Timezone:          "US/Pacific",
			InactiveTime:      2026,
			APIVersion:        "3.0",
			DeviceListVersion: 1,
			Fmly:              true,
		},
	}

	err := fmipService.api.Call(resource, map[string]string{
		"clientBuildNumber":     clientBuildNumber,
		"clientID":              clientID,
		"clientMasteringNumber": clientMasteringNumber,
		"dsid":                  dsid,
	}, payload, additionalHeaders)
	if err != nil {
		log.Fatalln(err)
	}
}

// {"serverContext":{"minCallbackIntervalInMS":5000,"enable2FAFamilyActions":false,"preferredLanguage":"en-gb","lastSessionExtensionTime":null,"enableMapStats":true,"callbackIntervalInMS":2000,"validRegion":true,"timezone":{"currentOffset":-25200000,"previousTransition":1520762399999,"previousOffset":-28800000,"tzCurrentName":"-07:00","tzName":"US/Pacific"},"authToken":null,"maxCallbackIntervalInMS":60000,"classicUser":false,"isHSA":true,"trackInfoCacheDurationInSecs":86400,"imageBaseUrl":"https://statici.icloud.com","minTrackLocThresholdInMts":100,"maxLocatingTime":90000,"sessionLifespan":900000,"info":"jJ+B16YUsHV0HRIUPcVnKQ7+n5CirsrMwU2hL7to504xj8DTqVQhbpn9p4QHSG71","prefsUpdateTime":1389369176042,"useAuthWidget":true,"clientId":"Y2xpZW50XzQxMTY0OTIxN18xNTI3MzQ1ODE5OTkw","enable2FAFamilyRemove":false,"serverTimestamp":1527345820092,"macCount":0,"deviceLoadStatus":"200","maxDeviceLoadTime":60000,"prsId":411649217,"showSllNow":false,"cloudUser":true,"enable2FAErase":false,"id":"server_ctx"},
//"clientContext":{"appName":"iCloud Find (Web)","appVersion":"2.0","timezone":"US/Pacific","inactiveTime":198,"apiVersion":"3.0","deviceListVersion":1,"fmly":true}}
// type serverContext struct {
//
// }

// type refreshClientRequestPayload struct {
// 	ServerContext serverContext `json:"serverContext"`
// 	ClientContext clientContext `json:"clientContext"`
// }

func (fmipService *FmipService) RefreshWebClient(xAppleCookiesHeader,
	clientBuildNumber, clientID, clientMasteringNumber, dsid string) {
	router := coreapi.NewRouter()
	router.DefaultRouter = coreapi.DefaultRouter
	router.RegisterFunc(200, fmipService.clientSuccess)
	resource := coreapi.NewResource(
		"/fmipservice/client/web/refreshClient?clientBuildNumber={{.clientBuildNumber}}&clientId={{.clientID}}&clientMasteringNumber={{.clientMasteringNumber}}&dsid={{.dsid}}",
		"POST", router)

	additionalHeaders := map[string]string{
		"Cookie": xAppleCookiesHeader,
	}

	payload := &initClientRequestPayload{
		ClientContext: clientContext{
			AppName:           "iCloud Find (Web)",
			AppVersion:        "2.0",
			Timezone:          "US/Pacific",
			InactiveTime:      2026,
			APIVersion:        "3.0",
			DeviceListVersion: 1,
			Fmly:              true,
		},
	}

	err := fmipService.api.Call(resource, map[string]string{
		"clientBuildNumber":     clientBuildNumber,
		"clientID":              clientID,
		"clientMasteringNumber": clientMasteringNumber,
		"dsid":                  dsid,
	}, payload, additionalHeaders)
	if err != nil {
		log.Fatalln(err)
	}
}
