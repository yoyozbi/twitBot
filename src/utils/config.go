package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)
type Track struct {
	Webhook 	string `json:"webhook"`
	Message 	string `json:"message"`
	Username 	string `json:"username"`
  WithReplies bool `json:"withReplies"`
  WithRetweets bool `json:"withRetweets"`
}
type Config struct {
	ConsumerKey    string `json:"twitterConsumerKey"`
	ConsumerSecret string `json:"twitterConsumerSecret"`
	AccessToken    string `json:"twitterToken"`
	AccessSecret   string `json:"twitterTokenSecret"`
	Track 		   []Track `json:"track"`
}

func (c Config) TwitterRule() string {
	var usernames []string
	for _, track := range c.Track {
		usernames = append(usernames, "from:" + track.Username)
	}

	return strings.Join(usernames, " OR ")
}


func LoadConfig() Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	var configData Config
	
	fileContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(fileContent, &configData)

	return configData
}
