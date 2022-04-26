/*
  Copyright (C) 2022 Zbinden Yohan 
*/

const { getConfig, sleep } = require("./utils");
const { sendMessageToWebhook } = require("./discord");
const { TwitterApi, ETwitterStreamEvent } = require("twitter-api-v2");

const config = getConfig();
async function load() {
  const client = new TwitterApi({
    appKey: config.twitterConsumerKey,
    appSecret: config.twitterConsumerSecret,
    accessToken: config.twitterToken,
    accessSecret: config.twitterTokenSecret,
  });
  const appClient = await client.appLogin();
  await editStreamRules(appClient, config);

  const stream = await appClient.v2.searchStream({ autoConnect: true });

  stream.on(ETwitterStreamEvent.Connected, () =>
    console.log("[TWITTER] connected")
  );

  stream.on(ETwitterStreamEvent.ReconnectAttempt, () => {
    console.log("[TWITTER] reconnecting");
  });

  stream.on(ETwitterStreamEvent.Data, async (eventData) => {
    if (eventData.data.id && eventData.data.text) {
      let tweet = await appClient.v2.singleTweet(eventData.data.id, {
        "tweet.fields": "author_id",
      });
      let author = await appClient.v2.user(tweet.data.author_id);
      let t = config.track.find((key) => key.username == author.data.username);
      if (t) {
        let link = `https://twitter.com/${author.data.username}/status/${eventData.data.id}`;
        let mess = t.message;
        mess = mess.replace("{link}", link);
        sendMessageToWebhook(t, mess);
      }
    } else {
      console.log(eventData);
    }
  });
  process.on("exit", () => {
    stream.close();
  });
}

async function editStreamRules(appClient, config) {
  const rule = config.track.map((key) => `from:${key.username}`).join(" OR ");
  const appliedRules = await appClient.v2.streamRules();
  if (appliedRules.data) {
    const ruleId = appliedRules.data.find((r) => r.value == rule);
    if (ruleId) {
      console.log(`[TWITTER] rule already exists`);
    } else {
      let ids = appliedRules.data.map((rule) => rule.id);
      await appClient.v2.updateStreamRules({
        delete: {
          ids: ids,
        },
      });
      await appClient.v2.updateStreamRules({
        add: [{ value: rule }],
      });
      console.log(
        `[TWITTER] another rule(s) existed removed it and added new rule`
      );
    }
  } else {
    await appClient.v2.updateStreamRules({
      add: [{ value: rule }],
    });
    console.log(`[TWITTER] rule added`);
    console.log(`[TWITTER] waiting 10s for rules to be applied`);
    await sleep(10000);
    console.log(`[TWITTER] rules applied`);
  }
}
module.exports = { load };
