const puppeteer = require("puppeteer");
const fs = require("fs");

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://lncn.org");
  const html = await page.content();
  fs.writeFileSync("test.txt", html);
  await browser.close();
})();
