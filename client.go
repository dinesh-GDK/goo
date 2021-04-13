package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func client() {

	var ip_addr string

	// check ip validity
	// fmt.Print("Enter IP:PORT of the host: ")
	// fmt.Scanf("%s", &ip_addr)
	ip_addr = "localhost:5000"

	addr, err := net.ResolveTCPAddr("tcp", ip_addr)
	error_handler(err)

	conn, err := net.DialTCP("tcp", nil, addr)
	error_handler(err)

	ban_status, err := bufio.NewReader(conn).ReadString('\n')
	error_handler(err)
	ban_status = strings.TrimSuffix(ban_status, "\n")

	if ban_status == ":b" {
		fmt.Println("-->> YOU HAVE BEEN BANNED FROM ### <<--")
		return
	}

	if !client_cipher_verification(conn) {
		return
	}

	user_name := client_set_user_name(conn)

	if user_name == ":f" {
		return
	}

	go client_receive_message(conn, user_name)
	client_send_message(conn, user_name)
}

func client_cipher_verification(conn *net.TCPConn) bool {

	var cipher string
	for {
		// fmt.Print("Enter cipher: ")
		// fmt.Scanf("%s", &cipher)
		cipher = "aaa"

		_, err := conn.Write([]byte(cipher + "\n"))
		error_handler(err)

		response, err := bufio.NewReader(conn).ReadString('\n')
		error_handler(err)
		response = strings.TrimSuffix(response, "\n")

		if response == ":m" {
			fmt.Println("-->> CIPHER MATCHED <<--")
			return true

		} else if response == ":le" {
			fmt.Println("-->> MAXIMUM LIMIT EXCEEDED FOR CIPHER VERIFICATION <<--")
			break

		} else {
			fmt.Println("-->> " + strings.Split(response, "#")[1] + " ATTEMPT(S) LEFT FOR CIPHER VERIFICATION <<--")
		}
	}

	return false
}

func client_set_user_name(conn *net.TCPConn) string {

	var user_name string
	for {
		fmt.Print("Enter User name: ")
		fmt.Scanf("%s", &user_name)

		_, err := conn.Write([]byte(user_name + "\n"))
		error_handler(err)

		response, err := bufio.NewReader(conn).ReadString('\n')
		error_handler(err)
		response = strings.TrimSuffix(response, "\n")

		if response == ":f" {
			fmt.Println("-->> Room Full <<--")
			return response

		} else if response == ":s" {
			break

		} else {
			fmt.Println("-->> User Name not available <<--")
		}
	}

	return user_name
}

func client_receive_message(conn *net.TCPConn, user_name string) {

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			clear_chat_line(user_name)
			fmt.Println("-->> DISCONNECTED <<--")
			os.Exit(99) // change
			return
		}

		clear_chat_line(user_name)
		fmt.Print(message)
		print_chat_line(user_name)
	}
}

func client_send_message(conn *net.TCPConn, user_name string) {

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

		if temp[0] == ':' {
			command_palette(temp)
			continue
		}

		_, err = conn.Write([]byte(string(text)))
		error_handler(err)
	}
}
