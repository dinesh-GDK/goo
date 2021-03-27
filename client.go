package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func client() {

	var err error
	var ip_addr string

	// check ip validity
	fmt.Print("Enter IP:PORT of the host: ")
	fmt.Scanf("%s", &ip_addr)

	addr, err := net.ResolveTCPAddr("tcp", ip_addr)
	error_handler(err)

	conn, err := net.DialTCP("tcp", nil, addr)
	error_handler(err)

	for {
		var user_name string
		fmt.Print("Enter User name: ")
		fmt.Scanf("%s", &user_name)

		_, err = conn.Write([]byte(user_name + "\n"))
		error_handler(err)

		response, err := bufio.NewReader(conn).ReadString('\n')
		error_handler(err)
		response = strings.TrimSuffix(response, "\n")

		if response == "set" {
			break
		} else {
			fmt.Println("Name already exists")
		}
	}

	// receive and print message
	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			error_handler(err)
			fmt.Print(message)
		}
	}()

	// send message
	for {

		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		error_handler(err)
		temp := strings.TrimSpace(string(text))

		_, err = conn.Write([]byte(string(text)))
		error_handler(err)

		if temp == ":stop" {
			conn.Close()
			os.Exit(1)
		}
	}

}
