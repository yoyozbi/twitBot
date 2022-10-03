package webhook

import (
	"encoding/json"
	"github.com/containrrr/shoutrrr"
	"github.com/yoyozbi/twitBot/src/twitterImpl"
	"github.com/yoyozbi/twitBot/src/utils"
	"log"
	"net/url"
	"strings"
)

type Body struct {
	Content  string `json:"content,omitempty"`
	IconURL  string `json:"icon_url,omitempty" `
	Username string `json:"username,omitempty"`
}

func Post(t utils.Track, tweet *twitterImpl.Tweet) error {
	msg := strings.ReplaceAll(t.Message, "{link}", tweet.GetUrl())
	url, err := url.ParseRequestURI(t.Url)
	if err != nil {
		return err
	}
	//we need to add a query parameter to send data as json
	//as explained here https://stackoverflow.com/questions/52170547/unable-to-set-query-parameters-manually-for-a-rest-api-using-mux
	//we cannot do it directly
	params := url.Query()
	params.Set("json", "Yes")
	url.RawQuery = params.Encode()

	body := Body{Content: msg}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = shoutrrr.Send(url.String(), string(b))
	if err != nil {
		return err
	}
	log.Println("[WEBHOOK] Sent tweet from @" + t.Username)

	return nil
}
