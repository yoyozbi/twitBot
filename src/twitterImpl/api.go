package twitterImpl

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/yoyozbi/twitBot/src/utils"
)

type TwitterApi struct {
	Config utils.Config
}

var accessToken string

const BASE_URL = "https://api.twitter.com"

//------------
//   tokens
// -----------
type Oauth2ClientResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

func RetrievAcessToken(apiKey string, apiSecretKey string) (Oauth2ClientResponse, error) {
	url := BASE_URL + "/oauth2/token"
	b := []byte("grant_type=client_credentials")
	req, err := http.NewRequest("post", url, bytes.NewBuffer(b))
	if err != nil {
		return Oauth2ClientResponse{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(apiKey+":"+apiSecretKey)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Oauth2ClientResponse{}, err
	}
	if resp.StatusCode == http.StatusForbidden {
		return Oauth2ClientResponse{}, errors.New("wrong tokens")
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Oauth2ClientResponse{}, err
	}
	var res Oauth2ClientResponse
	err = json.Unmarshal(responseData, &res)

	return res, err
}
