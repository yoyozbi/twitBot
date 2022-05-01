package twitterImpl

import (
	"context"
	"log"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/yoyozbi/twitBot/src/utils"
	"golang.org/x/oauth2/clientcredentials"
)

func Start(c chan *twitter.Tweet) {

	config := utils.LoadConfig()

	twitConfig := &clientcredentials.Config{
		ClientID:     config.ConsumerKey,
		ClientSecret: config.ConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := twitConfig.Client(context.Background())
	client := twitter.NewClient(httpClient)

	st := &Stream{
		client: httpClient,
	}

	_, err := st.SetRules([]Rule{{Value: config.TwitterRule()}})
	if err != nil {
		log.Fatal(err)
	}

	r := make(chan StreamResponse)
	go st.Connect(r)
	defer close(r)
	for {
		m := <-r
		iId, err := strconv.ParseInt(m.Data.Id, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		tweet, _, err := client.Statuses.Show(iId, &twitter.StatusShowParams{TrimUser: BoolPointer(false), IncludeEntities: BoolPointer(true)})

		c <- tweet
	}

}

func BoolPointer(b bool) *bool {
	return &b
}
