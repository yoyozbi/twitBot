/*
  Copyright (C) 2022 Zbinden Yohan 
*/
const { readFileSync } = require("fs");
const { join } = require("path");
function getConfig() {
  return JSON.parse(readFileSync(join(__dirname, "../config.json")));
}
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
module.exports = { getConfig, sleep };
