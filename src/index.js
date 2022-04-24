/*
  Copyright (C) 2022 Zbinden Yohan 
*/
const { getConfig } = require("./utils");
const tw = require("./twitter");

/*
-----------
| Twitter |
-----------
*/
(async () => {
  await tw.load();
})();
