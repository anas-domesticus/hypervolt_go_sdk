package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	authEndpoint = "https://kc.prod.hypervolt.co.uk/realms/retail-customers/protocol/openid-connect/token"
)

func GetToken(user, password string) (*Token, error) {
	postData := url.Values{}
	postData.Add("client_id", "home-assistant")
	postData.Add("grant_type", "password")
	postData.Add("scope", "openid profile email offline_access")
	postData.Add("username", user)
	postData.Add("password", password)

	resp, err := http.PostForm(authEndpoint, postData)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("got bad status at login: %d", resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	token := Token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
