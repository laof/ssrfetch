package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	var sw sync.WaitGroup
	var nodes []string
	sw.Add(2)

	go func() {
		fmt.Println("A")
		time.Sleep(5 * time.Second)
		nodes = append(nodes, "a")
		sw.Done()
	}()

	go func() {
		fmt.Println("b")
		time.Sleep(1 * time.Second)
		nodes = append(nodes, "b")
		sw.Done()
	}()

	sw.Wait()

	fmt.Println(nodes)
}

func main2() {

	var nodes []string

	req, err := httpRequest("https://www.zhi" + "mian" + "fei.com")
	if err != nil {
		return
	}
	defer req.Body.Close()
	html, err := io.ReadAll(req.Body)

	if err != nil {
		return
	}

	exp2 := regexp.MustCompile(`data-clipboard-text=".*?"`)
	str := exp2.FindAllString(string(html), -1)
	for _, v := range str {
		node := strings.Replace(v, `data-clipboard-text="`, "", 1)
		node = strings.Replace(node, `"`, "", 1)
		nodes = append(nodes, node)
	}

	fmt.Println(nodes)

}

func httpRequest(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// 设置请求頭
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	client := http.Client{}
	// Do sends an HTTP request and returns an HTTP response
	// 发起一个HTTP请求，返回一个HTTP响应
	return client.Do(request)
}
