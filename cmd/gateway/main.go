package main

import (
	tool "fetch"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

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

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {
		s := gets(r)
		w.Write([]byte(online(s)))
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

func parse(html, s string) string {

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

	return strings.Join(nodes, s)
}

func online(s string) string {
	res, err := http.Get(tool.Host)
	fmt.Println(tool.Host)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	str, _ := ioutil.ReadAll(res.Body)
	return parse(string(str), s)
}

func local(s string) string {
	txt, _ := ioutil.ReadFile("test.html")
	return parse(string(txt), s)
}
