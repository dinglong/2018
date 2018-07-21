package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read stdin error, %v\n", err)
		return
	}

	var buff bytes.Buffer
	if err := json.Indent(&buff, data, "", "    "); err != nil {
		fmt.Fprintf(os.Stderr, "json indent error, %v\n", err)
		return
	}

	fmt.Printf("%s\n", buff.String())
}
