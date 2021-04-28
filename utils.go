package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func error_handler(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(30)
	}
}

func get_ip() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("--> NOT CONNECTED TO INTERNET <--")
		return "", errors.New(":nc")
	}

	defer conn.Close()
	ip_address := conn.LocalAddr().(*net.UDPAddr).String()
	ip_address = strings.Split(ip_address, ":")[0]

	return ip_address, nil
}

func clear_chat_line(user_name string) {
	fmt.Print("\033[2K")
	fmt.Printf("\033[%dD", len(user_name)+7)
}

func print_chat_line(user_name string) {
	fmt.Printf("*%s -->> ", user_name)
}

func clear_chat_prev_line() {
	fmt.Print("\033[A\033[K")
}
