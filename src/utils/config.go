package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Track struct {
	Webhook      string `json:"webhook"`
	Url          string `json:"url"`
	Message      string `json:"message"`
	Username     string `json:"username"`
	WithReplies  bool   `json:"withReplies"`
	WithRetweets bool   `json:"withRetweets"`
}
type Config struct {
	ApiKey       string  `json:"twitterApiKey"`
	ApiKeySecret string  `json:"twitterApiKeySecret"`
	Track        []Track `json:"track"`
}

func (c Config) TwitterRule() string {
	var usernames []string
	for _, track := range c.Track {
		usernames = append(usernames, "from:"+track.Username)
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
	var fileContent []byte

	fileContent, err = ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(fileContent, &configData)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(configData.Track); i++ {
		if configData.Track[i].Webhook != "" && configData.Track[i].Url == "" {
			log.Println("[WARNING] Webhook has been deprecated in config.json please use Url instead")
			configData.Track[i].Url = configData.Track[i].Webhook
		} else if configData.Track[i].Webhook != "" && configData.Track[i].Url != "" {
			log.Println("[WARNING] Both Webhook and Url has been set in config.json only Url will be used")
		}
	}
	return configData
}
