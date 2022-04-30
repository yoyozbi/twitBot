package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/yoyozbi/twitBot/src/utils"
)
type WebhookBody struct {
	Content string `json:"content,omitempty"`
	IconURL string `json:"icon_url,omitempty" `
	Username string `json:"username,omitempty"`
}
func Post(t utils.Track, tweet *twitter.Tweet) error{
	msg := strings.ReplaceAll(t.Message,"{link}","https://twitter.com/" + t.Username + "/status/" + strconv.FormatInt(tweet.ID,10))
	body := WebhookBody{Content: msg}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Post(t.Webhook,"application/json", bytes.NewBuffer(b))
	if err != nil{
		return err
	}
	if resp.StatusCode != http.StatusNoContent{
		return errors.New("unexpected return code: " + resp.Status)
	}
	log.Println("[WEBHOOK] Sent tweet from @"+ t.Username)

	return nil
}