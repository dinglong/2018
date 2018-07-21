package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"unicode"

	"github.com/ghodss/yaml"
)

var format string

func init() {
	flag.StringVar(&format, "o", "json", "Format stdin to json or yaml")
	flag.Parse()
}

func main() {
	if format != "json" && format != "yaml" {
		fmt.Fprintf(os.Stderr, "error format %s, to use json or yaml\n", format)
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read stdin error, %v\n", err)
		return
	}

	if format == "json" {
		data, err := toJSON(data)
		if err != nil {
			fmt.Fprint(os.Stderr, "to json error, %v\n", err)
			return
		}

		var buff bytes.Buffer
		if err := json.Indent(&buff, data, "", "    "); err != nil {
			fmt.Fprintf(os.Stderr, "json indent error, %v\n", err)
			return
		}

		fmt.Printf("%s\n", buff.String())
	} else {
		data, err := toYAML(data)
		if err != nil {
			fmt.Fprint(os.Stderr, "to yaml error, %v\n", err)
			return
		}

		fmt.Printf("%s\n", string(data))
	}
}

// toJSON converts a single YAML document into a JSON document
// or returns an error. If the document appears to be JSON the
// YAML decoding path is not used (so that error messages are
// JSON specific).
func toJSON(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return data, nil
	}
	return yaml.YAMLToJSON(data)
}

func toYAML(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return yaml.JSONToYAML(data)
	}
	return data, nil
}

var jsonPrefix = []byte("{")

// hasJSONPrefix returns true if the provided buffer appears to start with
// a JSON open brace.
func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, jsonPrefix)
}

// Return true if the first non-whitespace bytes in buf is
// prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}
