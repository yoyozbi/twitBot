package twitterImpl

import (
	"encoding/json"
	"errors"
	"net/http"
)

const RULES_URL = BASE_URL + "/2/tweets/search/stream/rules"

//  types

type StreamRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

type StreamRules struct {
	Data []struct {
		Id    string `json:"id"`
		Value string `json:"value"`
	} `json:"data"`
	Meta struct {
		Sent        string `json:"sent"`
		ResultCount int    `json:"result_count"`
	} `json:"meta"`
}
type DeleteStreamRules struct {
	Ids []string `json:"ids,omitempty"`
}
type DeleteStreamRulesParams struct {
	Delete DeleteStreamRules `json:"delete,omitempty" binding:"required"`
}

type AddStreamRulesParams struct {
	Add []StreamRule `json:"add,omitempty" binding:"required"`
}

type UpdateStreamResponse struct {
	Meta struct {
		Sent    string `json:"sent"`
		Summary struct {
			Deleted    int `json:"deleted"`
			NotDeleted int `json:"not_deleted"`
			Created    int `json:"created"`
			NotCreated int `json:"not_created"`
			Valid      int `json:"valid"`
			NotValid   int `json:"not_valid"`
		} `json:"summary"`
	} `json:"meta"`
	Errors []struct {
		Value string `json:"value"`
		Id    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"errors"`
}

// methods

func (t *TwitterApi) GetAllStreamRules() (StreamRules, error) {
	_, responseData, err := t.makeHttpRequest("GET", RULES_URL, BodyParams{})
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

func (t *TwitterApi) UpdateStreamRule(params interface{}) (UpdateStreamResponse, error) {
	switch params.(type) {
	case DeleteStreamRulesParams:
	case AddStreamRulesParams:
	default:
		return UpdateStreamResponse{}, errors.New("unexpected param type")
	}
	body, err := json.Marshal(params)
	if err != nil {
		return UpdateStreamResponse{}, err
	}

	r, respBody, err := t.makeHttpRequest("POST", RULES_URL, BodyParams{ContentType: "application/json", Body: body})
	if r.StatusCode != http.StatusCreated && r.StatusCode != http.StatusOK {
		return UpdateStreamResponse{}, errors.New(string(respBody))
	}
	if err != nil {
		return UpdateStreamResponse{}, err
	}
	var res UpdateStreamResponse
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return UpdateStreamResponse{}, err
	}
	return res, nil
}

func (t *TwitterApi) SetStreamRules(r []StreamRule) (UpdateStreamResponse, error) {
	params := AddStreamRulesParams{
		Add: r,
	}
	return t.UpdateStreamRule(params)
}

func (t *TwitterApi) DeleteAllStreamRules() (UpdateStreamResponse, error) {
	rules, err := t.GetAllStreamRules()
	if err != nil {
		return UpdateStreamResponse{}, err
	}
	ids := make([]string, len(rules.Data))
	if len(ids) > 0 {
		for i, rule := range rules.Data {
			ids[i] = rule.Id
		}
		params := DeleteStreamRulesParams{
			Delete: DeleteStreamRules{Ids: ids},
		}
		return t.UpdateStreamRule(params)
	} else {
		return UpdateStreamResponse{}, nil
	}
}
