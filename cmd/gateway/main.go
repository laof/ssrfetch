package main

import (
	tool "fetch"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/apex/gateway"
)

func gets(r *http.Request) string {
	query := r.URL.Query()
	s := query.Get("s")

	if s == "" {
		return "\n"
	}

	return s
}

func main() {

	port := flag.Int("port", -1, "localhost port")

	flag.Parse()

	http.HandleFunc("/api/goto", func(w http.ResponseWriter, r *http.Request) {
		str := strings.Replace(r.RequestURI, "/api/goto?s=", "", 1)

		if str != "" {
			str, _ = url.QueryUnescape(str)
		}

		req, err := http.Get(str)
		if err != nil {
			w.Write([]byte("Get failed: " + str))
			return
		}
		defer req.Body.Close()
		data, err := io.ReadAll(req.Body)

		if err != nil {
			w.Write([]byte("ReadAll failed"))
			return
		}

		w.Write(data)

	})

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {
		s := gets(r)

		w.Write([]byte(node(s)))
	})

	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		hello := fmt.Sprintf("Hi%sI'm ok", gets(r))
		w.Write([]byte(hello))
	})

	if *port == -1 {
		fmt.Println("aws-lambda-go")
		gateway.ListenAndServe("n/a", nil)
	} else {
		sp := fmt.Sprintf(":%d", *port)
		fmt.Println("http://localhost" + sp)
		http.ListenAndServe(sp, nil)
	}

}

func node(s string) string {

	var sw sync.WaitGroup
	var nodes []string
	sw.Add(2)

	go func() {
		nodes = append(nodes, online()...)
		sw.Done()
	}()

	go func() {
		nodes = append(nodes, fei()...)
		sw.Done()
	}()

	sw.Wait()

	return strings.Join(nodes, s)

}

func parse(html string) []string {

	p := regexp.MustCompile(`<p>(.*)?</p>`)

	target := p.FindAllString(html, -1)

	// ssr
	ssr := regexp.MustCompile(`ssr://([^<]*)`)

	// ss
	ss := regexp.MustCompile(`ss://([^<]*)`)

	var data []string
	for _, p := range target {

		ssrtxt := ssr.FindAllString(p, -1)
		data = append(data, ssrtxt...)

		sstxt := ss.FindAllString(p, -1)
		data = append(data, sstxt...)

	}

	return data
}

func online() []string {
	res, err := httpRequest(tool.Host)
	fmt.Println(tool.Host)
	if err != nil {
		return make([]string, 0)
	}
	defer res.Body.Close()

	str, _ := io.ReadAll(res.Body)
	return parse(string(str))
}

// func local() []string {
// 	txt, _ := os.ReadFile("test.html")
// 	return parse(string(txt))
// }

func fei() []string {

	var data []string

	req, err := httpRequest("https://www.zhi" + "mian" + "fei.com")
	if err != nil {
		return data
	}
	defer req.Body.Close()
	html, err := io.ReadAll(req.Body)

	if err != nil {
		return data
	}

	exp2 := regexp.MustCompile(`data-clipboard-text=".*?"`)
	str := exp2.FindAllString(string(html), -1)
	for _, v := range str {
		ssr := strings.Replace(v, `data-clipboard-text="`, "", 1)
		ssr = strings.Replace(ssr, `"`, "", 1)
		data = append(data, ssr)
	}

	return data

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
