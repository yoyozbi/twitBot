# A simple app posting to a webhook when someone make a tweet

## How to use

1. Create a new app on Twitter from the [developer portal](https://developer.twitter.com/en/apps).
2. Copy the `Consumer Key`, `Consumer Secret`, `Access Token` and `Access Token scret` from the app.
3. Create a new interaction on a discord channel and copy the `Webhook URL` from the interaction.
4. Clone the repo
5. Create a config.json file a the root of the repo with the following content:

```json
{
  "twitterApiKey": "<YOUR_KEY>",
  "twitterApiKeySecret": "<YOUR_KEY>",
  "track": [
    {
      "webhook": "<YOUR_WEBHOOK_URL>",
      "message": "{link}",
      "username": "<THE_TWITTER_ACCOUNT>",
      "withReplies": false, //if you want to send tweet replies
      "withRetweets": false, //if you want to sent retweets
    }
  ]
}
```
