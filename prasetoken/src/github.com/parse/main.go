package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/parse/token"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read stdio error %v\n", err)
	}

	t, err := token.NewToken(string(data))
	if err != nil {
		log.Fatalf("new token error, %v\n", err)
	}

	h, err := json.Marshal(t.Header)
	if err != nil {
		log.Fatalf("marshal header error, %v\n", err)
	}
	printJson(h)

	c, err := json.Marshal(t.Claims)
	if err != nil {
		log.Fatalf("marshal claims error, %v\n", err)
	}
	printJson(c)
}

func printJson(data []byte) {
	var buff bytes.Buffer
	if err := json.Indent(&buff, data, "", "    "); err != nil {
		log.Fatalf("indent json error, %v\n", err)
	}
	fmt.Print(buff.String())
}
