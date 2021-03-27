package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func error_handler(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func get_ip() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	error_handler(err)

	defer conn.Close()
	ip_address := conn.LocalAddr().(*net.UDPAddr).String()
	ip_address = strings.Split(ip_address, ":")[0]

	return ip_address
}
