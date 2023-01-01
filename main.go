package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/apex/gateway/v2"
)

var port = ":7965"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.RequestURI()

		switch url {
		case "/api/get":
			w.Write([]byte(online()))
		case "/api/test":
			w.Write([]byte(local()))
		default:
			w.Write([]byte("hello world"))
		}

	})

	http.HandleFunc("api/get", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(online()))
	})

	http.HandleFunc("api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(local()))
	})

	fmt.Println("http://localhost" + port)
	gateway.ListenAndServe(port, nil)

}

func parse(html string) string {

	p := regexp.MustCompile(`<p>(.*)?</p>`)

	target := p.FindAllString(html, -1)

	// ssr
	ssr := regexp.MustCompile(`ssr://([^<]*)`)

	// ss
	ss := regexp.MustCompile(`ss://([^<]*)`)

	var nodes []string
	for _, p := range target {

		ssrtxt := ssr.FindAllString(p, -1)
		nodes = append(nodes, ssrtxt...)

		sstxt := ss.FindAllString(p, -1)
		nodes = append(nodes, sstxt...)

	}
	return strings.Join(nodes, ",")
}

func online() string {
	res, err := http.Get(host)
	fmt.Println(host)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	str, _ := ioutil.ReadAll(res.Body)
	return parse(string(str))
}

func local() string {
	txt, _ := ioutil.ReadFile("./test.html")
	return parse(string(txt))
}
