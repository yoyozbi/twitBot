package twitterImpl

import (
	"encoding/json"
	"errors"
  "log"
)


type User struct {
  Id       string `json:"id"`
  Name     string `json:"name"`
  Username string `json:"username"`
}

type UserResponse struct {
  Data User `json:"data"`
}

func (t *TwitterApi) GetUser(id string) (User, error) {
  url := BASE_URL + "/2/users/" + id
  resp, respBody, err := t.makeHttpRequest("get", url, BodyParams{})
  if err != nil {
    log.Println("salut")
    return User{}, err 
  }
  if resp.StatusCode == 404 {
    return User{}, errors.New("User not found")
  }

  var user UserResponse
  err = json.Unmarshal(respBody, &user)
  if err != nil {
    return User{}, err
  }
  return user.Data, nil
}

