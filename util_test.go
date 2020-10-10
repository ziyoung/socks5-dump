package main

import (
	"testing"
)

func TestParseAddr(t *testing.T) {
	inputs := []string{"http://baidu.com", "http://example.com:1234", "http://220.181.38.148:80"}
	for _, input := range inputs {
		addr, err := parseAddr(input)
		if err != nil {
			t.Fatal(err)
		}
		if len(addr) == 0 {
			t.Fatal("unable to parse address")
		}
	}
}
