package twitterImpl

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	STREAM_URL = BASE_URL + "/2/tweets/search/stream"
)

//-------
// types
//-------

type StreamResponse struct {
	Data struct {
		Id string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
	MatchingRules []struct {
		Id string `json:"id"`
		Tag string `json:"tag"`
	} `json:"matching_rules"`
}

//----------
// methods
//----------
func (t *TwitterApi) Connect(c chan StreamResponse) {
  req, err := http.NewRequest("get", STREAM_URL, nil)
  if err != nil {
    log.Fatal(err)
  }

  err = t.buildHeaders(req)
  if err != nil {
    log.Fatal(err)
  }
	resp, err := http.DefaultClient.Do(req) 
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)	
	}
	log.Println("[TWITTER] Stream Connected")
	dec := json.NewDecoder(resp.Body)
	for{
		var m StreamResponse
		err := dec.Decode(&m)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		c <- m
	}
}
