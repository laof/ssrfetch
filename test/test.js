const puppeteer = require("puppeteer");

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://lncn.org/api/ssr-list");
  const html = await page.content();
  fs.writeFileSync("test.txt", html);
  await browser.close();
})();
