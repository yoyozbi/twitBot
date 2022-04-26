/*
  Copyright (C) 2022 Zbinden Yohan 
*/
const { webLog } = require("./logger");
const { post } = require("axios");

async function sendMessageToWebhook(t, message) {
  try {
    await post(t.webhook, {
      content: message,
    });
  } catch (e) {
    webLog.error(e.response.data);
  }
  webLog.info(`sent new tweet from @${t.username}`);
}

module.exports = { sendMessageToWebhook };
