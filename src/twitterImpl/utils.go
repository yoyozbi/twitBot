package twitterImpl

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type BodyParams struct {
	ContentType string
	Body        []byte
}

//make a basic http request and return the body bytes
func (t *TwitterApi) makeHttpRequest(action string, url string, bodyParams BodyParams) (res *http.Response, resBody []byte, err error) {

	req, err := http.NewRequest(action, url, bytes.NewBuffer(bodyParams.Body))
	if err != nil {
		return nil, nil, err
	}
	err = t.buildHeaders(req)
	if err != nil {
		return nil, nil, err
	}
	if bodyParams.ContentType != "" {
		req.Header.Add("content_type", bodyParams.ContentType)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, responseBody, nil
}

//build headers with a new retrieved access token from the api or the one stored currently
func (t *TwitterApi) buildHeaders(req *http.Request) error {
	if accessToken == "" {
		t, err := RetrievAcessToken(t.Config.ConsumerKey, t.Config.ConsumerSecret)
		if err != nil {
			return err
		}
		accessToken = "Bearer " + t.AccessToken
	}
	req.Header.Add("Authorization", accessToken)
	return nil
}
