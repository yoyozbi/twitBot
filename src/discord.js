/*
  Copyright (C) 2022 Zbinden Yohan 
*/
const { post } = require("axios");

async function sendMessageToWebhook(t, message) {
  try {
    await post(t.webhook, {
      content: message,
    });
  } catch (e) {
    console.log(e.response.data);
  }
  console.log(`[WEBHOOK] sent new tweet from @${t.username}`);
}

module.exports = { sendMessageToWebhook };
