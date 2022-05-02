package twitterImpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	rulesURL = "https://api.twitter.com/2/tweets/search/stream/rules"
	streamURL = "https://api.twitter.com/2/tweets/search/stream"
)
type Stream struct {
	client *http.Client	
}

func (s *Stream) GetAllRules() StreamRules{
	resp, err := s.client.Get(rulesURL)	
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rules StreamRules
	err = json.Unmarshal(responseData, &rules)
	if err != nil {
		log.Fatal(err)
	}

	return rules
}

func (s *Stream)UpdateRule(params interface{}) (UpdateResponse,error) {
	switch params.(type) {
	case DeleteRulesParams:
	case AddRulesParams:
	default:
		return UpdateResponse{}, errors.New("unexpected param type")
	}

	body, err := json.Marshal(params)
	if err != nil {
		return UpdateResponse{}, err	
	}

	resp, err := s.client.Post(rulesURL,"application/json", bytes.NewBuffer(body))
	if err != nil {
		return UpdateResponse{},err
	}
	if resp.StatusCode != http.StatusCreated {
		return UpdateResponse{}, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UpdateResponse{}, err
	}

	var res UpdateResponse
	err = json.Unmarshal(responseData, &res)	
	if err != nil {	
		return UpdateResponse{}, err
	}	
	return res, nil
}

func (s *Stream) DeleteAllRules() (UpdateResponse, error) {
	rules := s.GetAllRules()	
	ids := make([]string, len(rules.Data))

	for i, rule := range rules.Data {
		ids[i] = rule.Id
	}
	params := DeleteRulesParams{
		Delete: DeleteRules{Ids:ids},
	}
	return s.UpdateRule(params)	

}



func (s *Stream) SetRules(r []Rule) (UpdateResponse, error) {
	//s.DeleteAllRules()
	params := AddRulesParams{
		Add: r,
	}
	return s.UpdateRule(params)	

}
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

func (s *Stream) Connect(c chan StreamResponse) {
	resp, err := s.client.Get(streamURL)
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
