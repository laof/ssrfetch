const puppeteer = require("puppeteer");
const { resolve, join } = require("path");
console.log(444444, __dirname);

console.log(66, __filename);

// 返回运行文件所在的目录
console.log("__dirname : " + __dirname);
// __dirname : /Desktop

// 当前命令所在的目录
console.log("resolve   : " + resolve("./"));
// resolve   : /workspace

// 当前命令所在的目录
console.log("cwd       : " + process.cwd());



(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://www.baidu.com");
  await page.pdf({ path: '/home/runner/work/ssrfetch/ssrfetch/test/te.pdf', format: "a4" });
  //   await page.pdf({ path: "/test/test3.pdf", format: "a4" });
  await browser.close();
})();
