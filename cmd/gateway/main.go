package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	tool "fetch"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/apex/gateway"
	"github.com/bitly/go-simplejson"
	"github.com/chromedp/chromedp"
	jwt "github.com/dgrijalva/jwt-go"
)

func gets(r *http.Request) string {
	query := r.URL.Query()
	s := query.Get("s")

	if s == "" {
		return "\n"
	}

	return s
}

func getPar(req *http.Request) map[string]string {
	decoder := json.NewDecoder(req.Body)
	var params map[string]string
	decoder.Decode(&params)
	return params
}

func main() {

	port := flag.Int("port", -1, "localhost port")

	flag.Parse()

	http.HandleFunc("/api/mp3/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		name := query.Get("name")
		list := search(name)
		a, _ := json.Marshal(list)
		w.Write(a)
	})
	// api/music?token=
	http.HandleFunc("/api/mp3/song", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		token := query.Get("token")
		w.Write([]byte(gomusic(token)))
	})

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

	http.HandleFunc("/api/signed", func(w http.ResponseWriter, r *http.Request) {
		params := getPar(r)
		data, err := createSigned(params["privateKey"], params["json"])
		if err == nil {
			w.Write([]byte(data))
		} else {
			w.Write([]byte(""))
		}
	})

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {
		s := gets(r)
		data := online()
		w.Write([]byte(strings.Join(data, s)))
	})

	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		hello := fmt.Sprintf("Hi%sI'm ok", gets(r))
		w.Write([]byte(hello))
	})

	http.HandleFunc("/api/lncn", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(lncn()))
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

// func node(s string) string {

// 	var sw sync.WaitGroup
// 	var nodes []string
// 	sw.Add(2)

// 	go func() {
// 		nodes = append(nodes, online()...)
// 		sw.Done()
// 	}()

// 	go func() {
// 		nodes = append(nodes, fei()...)
// 		sw.Done()
// 	}()

// 	sw.Wait()

// 	return strings.Join(nodes, s)

// }

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

func lncn() string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:])
	// chromedp.Flag("headless", true),
	// chromedp.Flag("disable-gpu", false),
	// chromedp.Flag("enable-automation", false),
	// chromedp.Flag("disable-extensions", false),

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string = "="

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://lncn.org/`),
		chromedp.OuterHTML(`body`, &res, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		return string(err.Error())
	}

	return res
}

type Music struct {
	Name  string `json:"name"`
	Dates string `json:"dates"`
	Act   string `json:"act"`
}

func dec(str string) string {
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&nbsp;", " ")
	return str
}

func search(name string) []Music {
	var list []Music

	res, err := http.Get("https://www.musicenc.com/?search=" + name)

	if err != nil {
		return list
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".list li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		name := s.Find("a").Text()
		name = dec(name)

		act := s.Find("span").Text()
		act = dec(act)

		dates, ok := s.Find("a").Attr("dates")

		if ok {
			m := Music{Name: name, Act: act, Dates: dates}
			list = append(list, m)
		}

	})

	return list

}

func gomusic(token string) string {
	res, err := http.Get("https://www.musicenc.com/searchr/?token=" + token)

	if err != nil {
		return ""
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return ""
	}

	var href string
	doc.Find(".downBu.secm3").Each(func(i int, s *goquery.Selection) {
		value, ok := s.Attr("href")
		if ok {
			href = value
		}
	})

	if strings.HasPrefix(href, "http") { // 我可以抱你吗
		return pics(href)
	}

	if strings.HasPrefix(href, "javascript:") { // 晴天
		return javascript(href)
	}

	return ""
}

// <script>var imgs="XJlcw==",lr="548323",pics="aHR0cHM6Ly9saW5rLm11c2ljZW5jLmNvbS8xNjMvP3NpZD0xNjc1ODAwMjI2JnRpbWU9MTY3NTU5MTc2NA==",sces="",domtitle="我可以抱你吗-MusicEnc",lrc="aHR0cHM6Ly93d3cubXVzaWNlbmMuY29tLzE2My8/cj0=",tl="";</script>
func pics(url string) string {
	var s = strconv.FormatInt(time.Now().UnixMilli(), 10)
	res, err := http.Get(url)

	if err != nil {
		return ""
	}

	defer res.Body.Close()

	doc, ok := io.ReadAll(res.Body)

	if ok != nil {
		return ""
	}

	str := string(doc)
	str = strings.ReplaceAll(str, "\"", s)
	str = strings.ReplaceAll(str, "'", s)

	a := strings.Split(str, "pics="+s)[1]
	value := strings.Split(a, s)[0]

	if value != "" {
		return download(value)
	}

	return ""
}

// 1,aHR0cHM6Ly9saW5rLm11c2ljZW5jLmNvbS8xNjMvP3NpZD0xNjc1ODAwMjI2JnRpbWU9MTY3NTU5MTc2NA==
// 2,https://link.musicenc.com/163/?sid=1675800226&time=1675591764
// 3,https://win-web-nf01-sycdn.kuwo.cn/6f14df20608f27fbd76a13f6f4aa6d16/63df8569/resource/n1/71/5/2472379884.mp3
func download(str string) string {
	url, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	res, err := http.Get(string(url))
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	mp3, e := io.ReadAll(res.Body)
	if e != nil {
		return ""
	}
	return string(mp3)
}

// 1,javascript:tps('aHR0cHM6Ly9hbnRpc2VydmVyLmt1d28uY24vYW50aS5zP2Zvcm1hdD1tcDN8YWFjJnJpZD0xMTg5ODAmYnI9MzIwa21wMyZ0eXBlPWNvbnZlcnRfdXJsJnJlc3BvbnNlPXJlcw==');
// 2,aHR0cHM6Ly9hbnRpc2VydmVyLmt1d28uY24vYW50aS5zP2Zvcm1hdD1tcDN8YWFjJnJpZD0xMTg5ODAmYnI9MzIwa21wMyZ0eXBlPWNvbnZlcnRfdXJsJnJlc3BvbnNlPXJlcw==
// 3,https://antiserver.kuwo.cn/anti.s?format=mp3|aac&rid=118980&br=320kmp3&type=convert_url&response=res
// 4,location url

func javascript(str string) string {
	str = strings.Split(strings.Split(str, "'")[1], "'")[0]
	url, _ := base64.StdEncoding.DecodeString(str)
	m := string(url)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(m)
	if err != nil {
		return ""
	}
	var l = res.Header.Get("Location")
	return l
}

/*
CreateSigned test
*/
func createSigned(key, jsonstr string) (tokenStr string, err error) {

	js, err := simplejson.NewJson([]byte(jsonstr))

	if err != nil {
		return "", err
	}

	maping, err := js.Map()

	if err != nil {
		return "", err
	}

	jsonmap := jwt.MapClaims{
		"timestamp": time.Now().Add(time.Hour*24).UnixNano() / 1e6,
	}

	for k := range maping {
		jsonmap[k] = maping[k]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jsonmap)

	token.Header = map[string]interface{}{
		"kid": "AWS",
		"alg": jwt.SigningMethodRS256.Alg(),
	}
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))

	return token.SignedString(privateKey)
}
