package gojsonparser_test

import (
	"fmt"
	"testing"

	gojsonparser "github.com/itsJakov/go-json-parser"
)

func TestAll(t *testing.T) {
	jsonStr := `[
		{
			"name": ["its", "Jakov"],
			"age": 12
		},
		{
			"name": "John\tDoe",
			"age": false,
			"friends": null
		}
	]`

	parsed, err := gojsonparser.ParseJson(jsonStr)
	if err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	fmt.Printf("%v\n", parsed)
}
