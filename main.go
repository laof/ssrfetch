package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/apex/gateway"
)

var port = flag.Int("port", -1, "specify a port")

func main() {
	fmt.Println(*port)

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(online()))
	})

	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am ok"))
	})

	if *port == -1 {
		fmt.Println("aws-lambda-go")
		gateway.ListenAndServe("n/a", nil)
	} else {
		sp := ":" + strconv.Itoa(*port)
		fmt.Println("http://localhost" + sp)
		http.ListenAndServe(sp, nil)
	}

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
	return strings.Join(nodes, "\n\n")
}

func online() string {
	res, err := http.Get(Host)
	fmt.Println(Host)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	str, _ := ioutil.ReadAll(res.Body)
	return parse(string(str))
}

func local() string {
	txt, _ := ioutil.ReadFile("test.html")
	return parse(string(txt))
}
