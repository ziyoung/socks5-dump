package main

import (
	"errors"
	"net"
	"net/url"
	"strconv"
)

func parseAddr(s string) ([]byte, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" {
		return nil, errors.New("only http is support")
	}

	host := u.Hostname()
	port := 80
	if u.Port() != "" {
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return nil, err
		}
	}

	var addr []byte
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			addr = make([]byte, 1+net.IPv4len+2)
			addr[0] = 1
			copy(addr[1:], ip4)
		} else {
			addr = make([]byte, 1+net.IPv6len+2)
			addr[0] = 4
			copy(addr[1:], ip)
		}
	} else {
		if len(host) > 255 {
			return nil, errors.New("host is too long")
		}
		addr = make([]byte, 1+1+len(host)+2)
		addr[0] = 3
		addr[1] = byte(len(host))
		copy(addr[2:], host)
	}

	addr[len(addr)-2], addr[len(addr)-1] = byte(port>>8), byte(port)

	return addr, nil
}
