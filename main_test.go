package main

import "testing"

func TestParseAddr(t *testing.T) {
	addr, err := parseAddr("http://baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(addr) == 0 {
		t.Fatal("unable to parse address")
	}
}
