package main

import (
	"encoding/base64"
	"strings"
)

var host string

func encode(code string) string {
	return base64.StdEncoding.EncodeToString([]byte(code))
}

func decode(code string) (str []byte, err error) {
	str, err = base64.StdEncoding.DecodeString(code)
	return
}

func init() {

	url := []string{
		"aHR0cHM6Ly9naXRodWIuY29tL0FsdmluOTk5O",
		"S9uZXctcGFjL3dpa2kvc3MlRTUlODUlOEQlRT",
		"glQjQlQjklRTglQjQlQTYlRTUlOEYlQjc=",
	}
	t, _ := decode(strings.Join(url, ""))
	host = string(t)
}
