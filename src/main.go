package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/yoyozbi/twitBot/src/twitterImpl"
	"github.com/yoyozbi/twitBot/src/utils"
	"github.com/yoyozbi/twitBot/src/webhook"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := utils.LoadConfig()

	c := make(chan *twitter.Tweet)
	go twitterImpl.Start(c)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		close(c)
		os.Exit(1)
	}()
	for {
		tweet := <-c
		for _, t := range config.Track {
			if t.Username == tweet.User.ScreenName {
				err := webhook.Post(t, tweet)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}
	/*signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	close(c)*/
}
