package twitterImpl

import (
	"bytes"
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
const BASE_URL = "https://api.twitter.com/";
//----------
// streamRules
//----------

const RULES_URL = BASE_URL + "/2/tweets/search/stream/rules"
//---------
//  types
//---------
type Rule struct {
	Value string `json:"value"`
	Tag string `json:"tag"`
}

type StreamRules struct {
	Data []struct {
		Id string `json:"id"`
		Value string `json:"value"`
	} `json:"data"`
	Meta struct {
		Sent string `json:"sent"`
		ResultCount int `json:"result_count"`
	} `json:"meta"`
}
type DeleteRules struct {	
	Ids []string `json:"ids,omitempty"`
}
type DeleteRulesParams struct {
	Delete DeleteRules `json:"delete,omitempty" binding:"required"`
}

type AddRulesParams struct {
	Add []Rule  `json:"add,omitempty" binding:"required"`
}

type UpdateResponse struct {
	Meta struct {
		Sent string `json:"sent"`
		Summary struct {
			Deleted int `json:"deleted"`
			NotDeleted int `json:"not_deleted"`
			Created int `json:"created"`
			NotCreated int `json:"not_created"`
			Valid int`json:"valid"`
			NotValid int `json:"not_valid"`
		} `json:"summary"`
	} `json:"meta"`
	Errors []struct {
		Value string `json:"value"`
		Id string `json:"id"`
		Title string `json:"title"`
		Type string `json:"type"`
	} `json:"errors"`
}
//---------
// methods
//---------

func (t *TwitterApi) GetAllStreamRules() (StreamRules, error)  {
  _,responseData, err := t.makeHttpRequest("GET", RULES_URL, BodyParams{})
  if err != nil {
    return StreamRules{}, err 
  }

  var rules StreamRules
  err = json.Unmarshal(responseData, &rules)
  if err != nil {
    return StreamRules{}, err
  }
  return rules, nil
}

func (t *TwitterApi) UpdateStreamRule(params interface{}) (UpdateResponse, error) {
  switch params.(type) {
    case DeleteRulesParams:
    case AddRulesParams:
    default:
      return UpdateResponse{}, errors.New("Unexpected param type")
    }
    body, err := json.Marshal(params)
  if err != nil {
    return UpdateResponse{}, err 
  }

  _, respBody, err := t.makeHttpRequest("POST", RULES_URL, BodyParams{contentType: "application/json", body: body})
  if err != nil {
    return UpdateResponse{}, err
  }

  var res UpdateResponse
  err = json.Unmarshal(respBody, &res)
  if err != nil {
    return UpdateResponse{}, err
  }
  return res, nil;  
}

func (t *StreamRules) SetStreamRules(r []Rule) (UpdateResponse, error) {
  params := AddRulesParams{
    Add: r,
  }
  return t.UpdateStreamRule(params)
}

func (t *TwitterApi) DeleteAllStreamRules() {
  
}


//--------
//  utils
//--------
type BodyParams struct {
  contentType string
  body        []byte
}
func (a *TwitterApi) makeHttpRequest(action string, url string, bodyParams BodyParams) (res *http.Response,resBody []byte, err error){

  req, err := http.NewRequest(action, url, bytes.NewBuffer(bodyParams.body))
  a.buildHeaders(req)
  if bodyParams.contentType != "" {
    req.Header.Add("content_type", bodyParams.contentType)
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
func (a *TwitterApi) buildHeaders(req *http.Request) error{
  if accessToken == "" {
    t, err := RetrievAcessToken(a.Config.ConsumerKey, a.Config.ConsumerSecret)
    if err != nil {
      return err
    }
    accessToken = t.TokenType + " " + t.AccessToken
  }
  req.Header.Add("Authorization", accessToken)
  return nil
}
//------------
//   tokens
// -----------
type Oauth2ClientResponse struct {
  TokenType string `json:"token_type"`
  AccessToken string `json:"access_token"`
}
func RetrievAcessToken(apiKey string, apiSecretKey string) (Oauth2ClientResponse, error){
  url := BASE_URL + "/oauth2/token";

  req, err := http.NewRequest("GET", url,nil )
  req.SetBasicAuth(apiKey,apiSecretKey)
  if err != nil {
    return Oauth2ClientResponse{},err
  }

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return Oauth2ClientResponse{},err
  }
  defer resp.Body.Close()
  responseData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
   return Oauth2ClientResponse{},err
  }
  var res Oauth2ClientResponse
  err = json.Unmarshal(responseData, &res)

  return res, err
}
