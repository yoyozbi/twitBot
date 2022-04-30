package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/yoyozbi/twitBot/src/twitterImpl"
)

func main() {
	c := make(chan *twitter.Tweet)
	go twitterImpl.Start(c)
	ch := make(chan os.Signal,1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	close(c)
}