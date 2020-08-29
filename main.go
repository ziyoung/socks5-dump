package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

var debugLog *log.Logger

func main() {
	var port int
	var url string
	var verbose bool
	flag.IntVar(&port, "port", 2048, "local proxy server port")
	flag.StringVar(&url, "url", "", "proxy url")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.Parse()
	if url == "" {
		log.Print("url is required")
		usage()
		return
	}

	debugLog = newDebugLog(verbose)
	address := "127.0.0.1:" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	err = handShake(conn, url)
	if err != nil {
		log.Fatal(err)
	}

	err = dialServer(conn, url)
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	log.Print("usage: socks5-dump [-port <port>] -url url [-verbose]")
}

func handShake(conn net.Conn, s string) error {
	// initial request
	log.Print("send initial request")
	initial := []byte{5, 1, 0}
	_, err := conn.Write(initial)
	if err != nil {
		return err
	}

	buf := make([]byte, 64)

	// initial response
	debugLog.Print("read initial response")
	_, err = io.ReadFull(conn, buf[:2])
	if err != nil {
		return err
	}
	if !bytes.Equal([]byte{5, 0}, buf[:2]) {
		return errors.New("invalid initial response")
	}

	// command request
	debugLog.Print("send command request")
	addr, err := parseAddr(s)
	if err != nil {
		return err
	}
	command := append([]byte{5, 1, 0}, addr...)
	_, err = conn.Write(command)
	if err != nil {
		return err
	}

	// command response
	log.Print("read command response")
	b := []byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0}
	_, err = io.ReadFull(conn, buf[:len(b)])
	if err != nil {
		return err
	}
	if !bytes.Equal(b, buf[:len(b)]) {
		return errors.New("invalid command response")
	}
	debugLog.Print("hand shake success")

	return nil
}

func dialServer(conn net.Conn, url string) error {
	debugLog.Print("send http request")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	err = req.Write(conn)
	if err != nil {
		return err
	}

	debugLog.Print("read http response")
	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, nil)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Print(hex.Dump(b))
	return nil
}
