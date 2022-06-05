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
	ApiKey    string `json:"twitterApiKey"`
	ApiKeySecret string `json:"twitterApiKeySecret"`
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
	err = json.Unmarshal(fileContent, &configData)
  if err != nil {
    log.Fatal(err)
  }
	return configData
}
