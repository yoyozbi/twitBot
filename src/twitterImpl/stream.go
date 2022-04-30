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