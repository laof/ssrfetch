package utils

import (
	"time"

	"github.com/bitly/go-simplejson"
	jwt "github.com/dgrijalva/jwt-go"
)

const key = `
test
`

type yamlData struct {
	Key   string
	Value string
}

func getPar(req *http.Request) map[string]string {
	decoder := json.NewDecoder(req.Body)
	var params map[string]string
	decoder.Decode(&params)
	return params
}
	params := getPar(req)
params["json"]
/*
CreateSigned test
*/
func CreateSigned(jsonstr string) (tokenStr string, err error) {

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
