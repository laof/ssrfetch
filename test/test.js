const puppeteer = require("puppeteer");


console.log(__dirname);

console.log(__filename);

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://www.baidu.com");
  await page.pdf({ path: "test/test.pdf", format: "a4" });
  //   await page.pdf({ path: "/test/test3.pdf", format: "a4" });
  await browser.close();
})();
