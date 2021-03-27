package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var clients []net.Conn
var user_names []string

func broadCast(message string) {

	fmt.Print(message)

	for i := range clients {
		_, err := clients[i].Write([]byte(message))
		error_handler(err)
	}
}

func handleConnection(conn *net.TCPConn) {

	var name string
	var err error

	fmt.Println(user_names)

	for {
		name, err = bufio.NewReader(conn).ReadString('\n')
		error_handler(err)

		name = strings.TrimSuffix(name, "\n")

		exist := false
		for _, user := range user_names {
			if name == user {
				exist = true
			}
		}

		if !exist {
			_, err = conn.Write([]byte(string("set\n")))
			error_handler(err)

			user_names = append(user_names, name)
			break

		} else {
			_, err = conn.Write([]byte(string("no_set\n")))
			error_handler(err)
		}

	}

	fmt.Println(user_names)

	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		if netData == ":stop" {
			fmt.Printf("Stopped serving %s\n", conn.RemoteAddr().String())
			break
		}

		broadCast(name + " >> " + netData)
	}
	conn.Close()
}

func host() {

	var PORT string

	// check for port validity
	fmt.Print("Enter the PORT to host: ")
	fmt.Scanf("%s", &PORT)

	addr, err := net.ResolveTCPAddr("tcp4", ":"+PORT)
	error_handler(err)

	listen, err := net.ListenTCP("tcp4", addr)
	error_handler(err)

	defer listen.Close()

	fmt.Println("Hosting at: " + get_ip() + PORT)

	var host_user_name string
	fmt.Print("Enter User name: ")
	fmt.Scanf("%s", &host_user_name)
	user_names = append(user_names, host_user_name)

	// send message
	go func() {
		for {

			text, err := bufio.NewReader(os.Stdin).ReadString('\n')
			error_handler(err)
			temp := strings.TrimSpace(string(text))

			broadCast(host_user_name + " >> " + text)

			if temp == ":stop" {
				listen.Close()
				os.Exit(1)
			}
		}
	}()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			return
		}
		clients = append(clients, conn)

		go handleConnection(conn)
	}
}
