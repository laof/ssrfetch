const puppeteer = require("puppeteer");

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://www.baidu.com");
  await page.pdf({ path: "test.pdf", format: "a4" });

  await browser.close();
})();
