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
)

func handShake(conn net.Conn, s string) error {
	// initial request
	initial := []byte{5, 1, 0}
	_, err := conn.Write(initial)
	if err != nil {
		return err
	}

	buf := make([]byte, 64)

	// initial response
	_, err = io.ReadFull(conn, buf[:2])
	if err != nil {
		return err
	}
	if !bytes.Equal([]byte{5, 0}, buf[:2]) {
		return errors.New("invalid initial response")
	}

	// command request
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
	b := []byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0}
	_, err = io.ReadFull(conn, buf[:len(b)])
	if err != nil {
		return err
	}
	if !bytes.Equal(b, buf[:len(b)]) {
		return errors.New("invalid command response")
	}

	return nil
}

func dialServer(conn net.Conn, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	err = req.Write(conn)
	if err != nil {
		return err
	}

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

func main() {
	var port int
	var url string
	flag.IntVar(&port, "port", 2048, "local proxy server port")
	flag.StringVar(&url, "url", "", "proxy url")
	flag.Parse()
	if url == "" {
		log.Fatalf("url is required\nusage: socks5-dump [-port <port>] -url url")
	}

	address := "127.0.0.1:" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)
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
