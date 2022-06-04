package twitterImpl
import (
  "errors"
  "encoding/json"
)

type ReferencedTweets struct {
  Type string `json:"type" binding:"required"`
  Id   string `json:"id" binding:"required"`
}
type Tweet struct{
  Id   string `json:"id" binding:"required"` 
  Text string `json:"text" binding:"required"`
  AuthorId  string `json:"author_id"`  
  ConversationId string `json:"conversation_id"`
  InReplyToUserId string `json:"in_reply_to_user_id"`
  ReferencedTweets []ReferencedTweets `json:"referenced_tweets"`
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
  var tweet Tweet
  err = json.Unmarshal(respBody, &tweet)
  if err != nil {
    return Tweet{}, err
  }
  return tweet, nil
}
