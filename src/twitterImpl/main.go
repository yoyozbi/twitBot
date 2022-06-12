package twitterImpl

import (
	"github.com/yoyozbi/twitBot/src/utils"
	"log"
)

func Start(c chan *Tweet) {

	config := utils.LoadConfig()

	t := TwitterApi{
		Config: config,
	}
	_, err := t.DeleteAllStreamRules()
	if err != nil {
		log.Fatal(err)
	}
	_, err = t.SetStreamRules([]StreamRule{{Value: config.TwitterRule()}})
	if err != nil {
		log.Fatal(err)
	}

	r := make(chan StreamResponse)
	go t.Connect(r)
	defer close(r)
	for {
		m := <-r

		tweet, err := t.GetTweet(m.Data.Id)
		if err != nil {
			log.Fatal(err)
		}
		c <- &tweet
	}

}
