package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ChimeraCoder/gojson"
)

func TestSimpleJson(t *testing.T) {
	var b []byte
	data, _ := ioutil.ReadFile("./test.json")
	i := strings.NewReader(string(data))
	b, err := gojson.Generate(i, gojson.ParseJson, "TestStruct", "main", []string{"bson", "json"}, true, true)
	if err != nil {
		t.Error("Generate() error:", err)
	}

	ioutil.WriteFile("gojsonStruct.go", b, 0644)

}
