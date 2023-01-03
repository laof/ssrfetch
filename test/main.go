package main

import (
	"fmt"
	"net/url"
)

var (
	a = "https://laof.github.io/flutter.txt"
)

func main() {
	enEscapeUrl, _ := url.QueryUnescape(a)
	fmt.Println("解码:", enEscapeUrl)
}
