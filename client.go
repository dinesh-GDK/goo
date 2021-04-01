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
	// fmt.Print("Enter IP:PORT of the host: ")
	// fmt.Scanf("%s", &ip_addr)
	ip_addr = "localhost:5000"

	addr, err := net.ResolveTCPAddr("tcp", ip_addr)
	error_handler(err)

	conn, err := net.DialTCP("tcp", nil, addr)
	error_handler(err)

	var user_name string
	for {

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
			fmt.Println("User name not available")
		}
	}

	// receive and print message
	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				clear_chat_line(user_name)
				fmt.Println("-->> DISCONNECTED <<--")
				os.Exit(34)
			}

			clear_chat_line(user_name)
			fmt.Print(message)
			print_chat_line(user_name)
		}
	}()

	// send message
	for {

		print_chat_line(user_name)

		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		error_handler(err)
		temp := strings.TrimSpace(string(text))

		clear_chat_prev_line()

		fmt.Print(user_name + " >> " + text)

		if temp == ":stop" {
			conn.Close()
			return
		}

		_, err = conn.Write([]byte(string(text)))
		error_handler(err)
	}

}
