package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	fmt.Println(online())
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
