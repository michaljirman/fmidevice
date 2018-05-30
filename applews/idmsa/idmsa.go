package idmsa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/michaljirman/coreapi"
)

type IdmsaService struct {
	api                    *coreapi.API
	XappleSessionToken     string
	XappleIDSessionID      string
	XappleIDAccountCountry string
}

func NewIdmsaService(baseURL string) (*IdmsaService, error) {
	// https://idmsa.apple.com/appleauth/auth/signin?widgetKey=83545bf919730e51dbfba2aaaaaaaaaa
	return &IdmsaService{
		api: NewIdmsaAPI(baseURL),
	}, nil
}

type accountLoginUIRequestPayload struct {
	AccountName string   `json:"accountName"`
	RememberMe  bool     `json:"rememberMe"`
	Password    string   `json:"password"`
	TrustTokens []string `json:"trustTokens"`
}

// serviceErrorsResponse struct
// {
// 	"serviceErrors": [{
// 		"code": "-20101",
// 		"message": "Your Apple ID or password was incorrect."
// 	}]
// }
type serviceErrorsResponse struct {
	ServiceErrors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"serviceErrors"`
}

func (idmsaService *IdmsaService) accountLoginUISuccess(resp *http.Response) error {
	if sessionToken := resp.Header.Get("X-Apple-Session-Token"); sessionToken != "" {
		idmsaService.XappleSessionToken = sessionToken
	} else {
		return fmt.Errorf("Unable to retrieve X-Apple-Session-Token")
	}

	if sessionID := resp.Header.Get("X-Apple-ID-Session-Id"); sessionID != "" {
		idmsaService.XappleIDSessionID = sessionID
	} else {
		return fmt.Errorf("Unable to retrieve X-Apple-Session-Id")
	}

	if accountCountry := resp.Header.Get("X-Apple-ID-Account-Country"); accountCountry != "" {
		idmsaService.XappleIDAccountCountry = accountCountry
	} else {
		return fmt.Errorf("Unable to retrive X-Apple-ID-Account-Country")
	}
	return nil
}

//2018/05/28 09:17:52 Unauthorized [Code: -20101, Message: Your Apple ID or password was incorrect.]
func (idmsaService *IdmsaService) accountLoginUIUnauthorized(resp *http.Response) error {
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respContent := serviceErrorsResponse{}
	json.Unmarshal(content, &respContent)

	var serviceErrors []string
	for _, err := range respContent.ServiceErrors {
		serviceErrors = append(serviceErrors, fmt.Sprintf("Code: %s, Message: %s", err.Code, err.Message))
	}
	return fmt.Errorf("Unauthorized %s", serviceErrors)
}

func (idmsaService *IdmsaService) AccountLoginUI(accountLoginUIResource, accountName, password string) {
	trustTokens := []string{}
	payload := &accountLoginUIRequestPayload{
		AccountName: accountName,
		RememberMe:  false,
		Password:    password,
		TrustTokens: trustTokens,
	}
	router := coreapi.NewRouter()
	router.DefaultRouter = coreapi.DefaultRouter
	router.RegisterFunc(401, idmsaService.accountLoginUIUnauthorized)
	router.RegisterFunc(409, idmsaService.accountLoginUISuccess)
	// https://idmsa.apple.com/appleauth/auth/signin?widgetKey=83545bf919730e51dbfba2aaaaaaaaaa
	resource := coreapi.NewResource(accountLoginUIResource, "POST", router)

	err := idmsaService.api.Call(resource, nil, payload, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
