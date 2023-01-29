const puppeteer = require("puppeteer");
const fs = require("fs");

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://lncn.org");
  await page.waitForSelector(".ssr-btn-bar button");
  const html = await page.content();
  fs.writeFileSync("test.txt", html);
  await browser.close();
})();
