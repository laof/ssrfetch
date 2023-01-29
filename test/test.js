const puppeteer = require("puppeteer");
const fs = require("fs");

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://lncn.org", {
    timeout: 1000 * 30,
    waitUntil: "networkidle0",
  });

  await page.reload({
    timeout: 1000 * 30,
    waitUntil: "networkidle0",
  });

  //   await page.waitForSelector(".ssr-btn-bar button");
  //   const html = await page.content();
  //   fs.writeFileSync("test.txt", html);
  await page.pdf({ path: "test.pdf" });
  await browser.close();
})();
