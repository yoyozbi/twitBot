# A simple app posting to a webhook when someone make a tweet## Self-hosting

### Standalone

1. Clone the repo (`git clone https://github.com/yoyozbi/twitbot.git`)
2. Create a new app on Twitter from the [developer portal](https://developer.twitter.com/en/apps).
3. Create a new interaction on a discord channel and copy the `Webhook URL` from the interaction. 
4. Create a config.json file at the root of the repo with the following content (notification uses [shoutrr](https://containrrr.dev/shoutrrr/v0.5/services/overview/) and should work with every services but only tested with discord)

```json
{
  "twitterApiKey": "<YOUR_KEY>",
  "twitterApiKeySecret": "<YOUR_KEY>",
  "track": [
    {
      "url": "<SHOUTRRR_URL>",
      "message": "{link}",
      "username": "<THE_TWITTER_ACCOUNT>",
      "withReplies": false, //if you want to send tweet replies
      "withRetweets": false, //if you want to sent retweets
    }
  ]
}
```
5. Build project with `go build src/main.go`
6. Run it with `chmod +x ./main && ./main` on linux or `./main.exe` in Windows

### Docker

1. Create the `config.json` file by following the standalone installation from step 2
2. Create a `docker-compose.yml` file and map the `config.json` file:

```yaml
version: '3'
services:
  image: ghcr.io/yoyozbi/twitbot:latest
  container_name: twitbot
  restart: unless-stopped
  volumes:
    - ./config.json:/data/config.json
```
3. run the app with `docker-compose up -d`