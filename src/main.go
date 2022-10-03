package main

import (
	"github.com/yoyozbi/twitBot/src/twitterImpl"
	"github.com/yoyozbi/twitBot/src/utils"
	"github.com/yoyozbi/twitBot/src/webhook"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := utils.LoadConfig()

	c := make(chan *twitterImpl.Tweet)
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
			if t.Username == tweet.Author.Username {
				if !t.WithReplies && tweet.IsReplied() {
					continue
				}
				if !t.WithRetweets && tweet.IsRetweeted() {
					continue
				}
				webhook.Post(t, tweet)
			}
		}
	}
}
