package coreapi

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
)

// AuthToken for token authentication
type AuthToken struct {
	Token string
}

// AuthBasic for basic authentication
type AuthBasic struct {
	Username string
	Password string
}

// Authentication interface
type Authentication interface {
	AuthorizationHeader() string // "basic <base64-encoded string>"
}

// NewAuthToken creates a new auth token struct
func NewAuthToken(token string) *AuthToken {
	return &AuthToken{
		Token: token,
	}
}

// AuthorizationHeader returns the token header value
func (a *AuthToken) AuthorizationHeader() string {
	return fmt.Sprintf("Token %s", a.Token)
}

// NewAuthBasic creates a new auth basic struct
func NewAuthBasic(username, password string) *AuthBasic {
	return &AuthBasic{
		Username: username,
		Password: password,
	}
}

// AuthorizationHeader returns the basic auth header value
func (a *AuthBasic) AuthorizationHeader() string {
	buffer := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, buffer)
	encContent := fmt.Sprintf("%s:%s", a.Username, a.Password)
	enc.Write([]byte(encContent))
	enc.Close()
	content, err := ioutil.ReadAll(buffer)
	if err != nil {
		log.Fatalln("Read failed:", err)
	}
	return fmt.Sprintf("Basic %s", string(content))
}
