package main

import (
	"io/ioutil"
	"log"
	"os"
)

func newDebugLog(debug bool) *log.Logger {
	wr := ioutil.Discard
	if debug {
		wr = os.Stderr
	}
	return log.New(wr, "", log.LstdFlags)
}
