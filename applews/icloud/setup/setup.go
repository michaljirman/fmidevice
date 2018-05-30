package setup

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/michaljirman/coreapi"
)

type SetupService struct {
	api                   *coreapi.API
	SetupRequestResponse  *setupRequestResponse
	AccountLoginResponse  *accountLoginResponse
	clientBuildNumber     string
	clientID              string
	clientMasteringNumber string
	xAppleCookies         []string
}

func NewSetupService(baseURL, clientBuildNumber, clientID, clientMasteringNumber string) (*SetupService, error) {
	return &SetupService{
		SetupRequestResponse:  nil,
		AccountLoginResponse:  nil,
		api:                   NewSetupAPI(baseURL),
		clientBuildNumber:     clientBuildNumber,
		clientID:              clientID,
		clientMasteringNumber: clientMasteringNumber,
	}, nil
}

type setupRequestResponse struct {
	Success     bool `json:"success"`
	RequestInfo []struct {
		Country         string `json:"country"`
		TimeZone        string `json:"timeZone"`
		IsAppleInternal bool   `json:"isAppleInternal"`
		Region          string `json:"region"`
	} `json:"requestInfo"`
	ConfigBag struct {
		Urls struct {
			AccountCreateUI     string `json:"accountCreateUI"`
			AccountLoginUI      string `json:"accountLoginUI"`
			AccountLogin        string `json:"accountLogin"`
			AccountRepairUI     string `json:"accountRepairUI"`
			DownloadICloudTerms string `json:"downloadICloudTerms"`
			RepairDone          string `json:"repairDone"`
			VettingURLForEmail  string `json:"vettingUrlForEmail"`
			AccountCreate       string `json:"accountCreate"`
			GetICloudTerms      string `json:"getICloudTerms"`
			VettingURLForPhone  string `json:"vettingUrlForPhone"`
		} `json:"urls"`
		AccountCreateEnabled bool `json:"accountCreateEnabled"`
	} `json:"configBag"`
	Error string `json:"error"`
}

func (ss *SetupService) GetAccountLoginUIHostAndResource() (string, string, error) {
	loginUrl := ss.SetupRequestResponse.ConfigBag.Urls.AccountLoginUI
	u, err := url.Parse(loginUrl)
	if err != nil {
		return "", "", err
	}
	baseURL := u.Scheme + "://" + u.Host
	resource := strings.Split(loginUrl, baseURL)[1]
	return baseURL, resource, nil
}

func (ss *SetupService) GetAccountLoginHostAndResource() (string, string, error) {
	loginUrl := ss.SetupRequestResponse.ConfigBag.Urls.AccountLogin
	u, err := url.Parse(loginUrl)
	if err != nil {
		return "", "", err
	}
	baseURL := u.Scheme + "://" + u.Host
	resource := strings.Split(loginUrl, baseURL)[1]
	return baseURL, resource, nil
}
func (ss *SetupService) GetXappleCookiesHeader() string {
	var xAppleCookiesHeader string
	for _, xAppleCookie := range ss.xAppleCookies {
		xAppleCookiesHeader += xAppleCookie + ";"
	}
	return xAppleCookiesHeader
}

func (ss *SetupService) validateSuccess(resp *http.Response) error {
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	ss.SetupRequestResponse = &setupRequestResponse{}
	json.Unmarshal(content, ss.SetupRequestResponse)
	return nil
}

////https://setup.icloud.com/setup/ws/1/validate?clientBuildNumber=1809Project50&clientId=A389681A-49C8-4436-8721-4B296F6B95D9&clientMasteringNumber=1809B29
func (ss *SetupService) Validate() {
	router := coreapi.NewRouter()
	router.DefaultRouter = coreapi.DefaultRouter
	router.RegisterFunc(421, ss.validateSuccess)
	resource := coreapi.NewResource(
		"/setup/ws/1/validate?clientBuildNumber={{.clientBuildNumber}}&clientId={{.clientID}}&clientMasteringNumber={{.clientMasteringNumber}}", "POST", router)

	err := ss.api.Call(resource, map[string]string{
		"clientBuildNumber":     ss.clientBuildNumber,
		"clientID":              ss.clientID,
		"clientMasteringNumber": ss.clientMasteringNumber,
	}, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

type accountLoginRequestPayload struct {
	DsWebAuthToken     string `json:"dsWebAuthToken"`
	AccountCountryCode string `json:"accountCountryCode"`
}

func (ss *SetupService) accountLoginSuccess(resp *http.Response) error {
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ss.AccountLoginResponse = &accountLoginResponse{}
	json.Unmarshal(content, ss.AccountLoginResponse)

	for _, setCookie := range resp.Header["Set-Cookie"] {
		ss.xAppleCookies = append(ss.xAppleCookies, setCookie)
	}
	return nil
}

type accountLoginResponse struct {
	DsInfo struct {
		Dsid string `json:"dsid"`
	} `json:"dsInfo"`
	Webservices struct {
		Findme struct {
			URL    string `json:"url"`
			Status string `json:"status"`
		} `json:"findme"`
	} `json:"webservices"`
}

func (ss *SetupService) AccountLogin(accountLoginUrlResource, xAppleSessionToken, xAppleIDAccountCountry string) {
	payload := &accountLoginRequestPayload{
		DsWebAuthToken:     xAppleSessionToken,
		AccountCountryCode: xAppleIDAccountCountry,
	}
	router := coreapi.NewRouter()
	router.DefaultRouter = coreapi.DefaultRouter
	router.RegisterFunc(200, ss.accountLoginSuccess)
	resource := coreapi.NewResource(
		accountLoginUrlResource+"?clientBuildNumber={{.clientBuildNumber}}&clientId={{.clientID}}&clientMasteringNumber={{.clientMasteringNumber}}", "POST", router)

	err := ss.api.Call(resource, map[string]string{
		"clientBuildNumber":     ss.clientBuildNumber,
		"clientID":              ss.clientID,
		"clientMasteringNumber": ss.clientMasteringNumber,
	}, payload, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
