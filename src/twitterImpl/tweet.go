package twitterImpl
import (
  "errors"
  "encoding/json"
  "log"
)

type ReferencedTweets struct {
  Type string `json:"type" binding:"required"`
  Id   string `json:"id" binding:"required"`
}
type TweetRaw struct{
  Id               string `json:"id" binding:"required"` 
  Text             string `json:"text" binding:"required"`
  AuthorId         string `json:"author_id"`  
  ConversationId   string `json:"conversation_id"`
  InReplyToUserId  string `json:"in_reply_to_user_id"`
  ReferencedTweets []ReferencedTweets `json:"referenced_tweets"`
}
type Tweet struct {
  Id               string
  Text             string
  Author           User
  ConversationId   string
  InReplyToUser    User
  ReferencedTweets []ReferencedTweets
}
type TweetResponse struct {
  Data TweetRaw `json:"data"`
}
func (t *Tweet) GetUrl() string {
  return ("https://twitter.com/" + t.Author.Username + "/status/" + t.Id)
}
func (t *Tweet) IsRetweeted() bool { 
  for _, ref := range t.ReferencedTweets { 
    if ref.Type == "retweeted" { 
      return true
    }
  }
  return false
}
func (t *Tweet) IsReplied() bool { 
  for _, ref := range t.ReferencedTweets { 
    if ref.Type == "replied_to" { 
      return true
    }
  }
  return false
}

func (t *TwitterApi) GetTweet(id string) (Tweet, error){
  url := BASE_URL + "/2/tweets/" + id + "?expansions=author_id,referenced_tweets.id,in_reply_to_user_id"
  resp, respBody, err := t.makeHttpRequest("get", url, BodyParams{})
  if err != nil {
    return Tweet{}, err
  }
  if resp.StatusCode == 404 {
    return Tweet{}, errors.New("Tweet not found")
  }
  var rawTweet TweetResponse
  err = json.Unmarshal(respBody, &rawTweet)
  if err != nil {
    log.Println(respBody)
    return Tweet{}, err
  }
  var author User
  if rawTweet.Data.AuthorId != "" {
    author, err = t.GetUser(rawTweet.Data.AuthorId)
    if err != nil {
      return Tweet{Id: rawTweet.Data.Id, Text: rawTweet.Data.Text, ConversationId: rawTweet.Data.ConversationId, ReferencedTweets: rawTweet.Data.ReferencedTweets}, err
    }
  }
  var InReplyToUser User
  if rawTweet.Data.InReplyToUserId != "" {
    InReplyToUser, err = t.GetUser(rawTweet.Data.InReplyToUserId)
    if err != nil {
      return Tweet{Id: rawTweet.Data.Id, Text: rawTweet.Data.Text, ConversationId: rawTweet.Data.ConversationId, ReferencedTweets: rawTweet.Data.ReferencedTweets}, err
    }
  }
  return Tweet{Id: rawTweet.Data.Id, Text: rawTweet.Data.Text, Author: author, ConversationId: rawTweet.Data.ConversationId, InReplyToUser: InReplyToUser, ReferencedTweets: rawTweet.Data.ReferencedTweets}, nil 
}
