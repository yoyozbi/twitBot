/*
  Copyright (C) 2022 Zbinden Yohan 
*/
const { readFileSync } = require("fs");
const { join } = require("path");
function getConfig() {
  return JSON.parse(readFileSync(join(__dirname, "../config.json")));
}

module.exports = { getConfig };
