package twitterImpl

import (
  "log"
	"github.com/yoyozbi/twitBot/src/utils"
)

func Start(c chan *Tweet) {

	config := utils.LoadConfig()

  //oAuth1config := oauth1.NewConfig()
	//twitConfig := &clientcredentials.Config{
	//	ClientID:     config.ConsumerKey,
	//	ClientSecret: config.ConsumerSecret,
	//	TokenURL:     "https://api.twitter.com/oauth2/token",
	//}
	//client := twitter.NewClient(httpClient)
  t := TwitterApi{
    Config: config,
  }	

	_, err := t.SetStreamRules([]StreamRule{{Value: config.TwitterRule()}})
	if err != nil {
		log.Fatal(err)
	}

	r := make(chan StreamResponse)
	go t.Connect(r)
	defer close(r)
	for {
		m := <-r

    tweet, err := t.GetTweet(m.Data.Id)
    if err != nil{
      log.Fatal(err)
    }
    c <- &tweet
	}

}

func BoolPointer(b bool) *bool {
	return &b
}
