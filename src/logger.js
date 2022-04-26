const { join } = require("path");
const bunyan = require("bunyan");

const logPath = join(__dirname, "../logs.log");

const streams = [
  {
    level: "info",
    stream: process.stdout,
  },
  {
    level: "info",
    path: logPath,
  },
];

const twitLog = bunyan.createLogger({ name: "twitter", streams: streams });
const webLog = bunyan.createLogger({ name: "webhook", streams: streams });
const othLog = bunyan.createLogger({ name: "other", streams: streams });
module.exports = { twitLog, webLog, othLog };
