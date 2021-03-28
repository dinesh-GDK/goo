package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var clients []*net.TCPConn
var user_names []string
var host_user_name string

func broadCast(message string, exclude_user string) {

	fmt.Print(message)

	for i := range clients {

		if exclude_user == user_names[i+1] {
			continue
		}

		_, err := clients[i].Write([]byte(message))
		error_handler(err)
	}
}

func handleConnection(conn *net.TCPConn) {

	var name string
	var err error

	// set user name
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

	// read from client
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

		fmt.Print("\033[2K")
		fmt.Printf("\033[%dD", len(host_user_name)+6)
		broadCast(name+" >> "+netData, name)
		fmt.Printf("%s -->> ", host_user_name)
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

	fmt.Println("Hosting at: " + get_ip() + ":" + PORT)

	fmt.Print("Enter User name: ")
	fmt.Scanf("%s", &host_user_name)
	user_names = append(user_names, host_user_name)

	// send message
	go func() {
		for {

			fmt.Printf("%s -->> ", host_user_name)

			text, err := bufio.NewReader(os.Stdin).ReadString('\n')
			error_handler(err)
			temp := strings.TrimSpace(string(text))

			fmt.Print("\033[A")
			fmt.Printf("\033[K")

			broadCast(host_user_name+" >> "+text, "")

			if temp == ":stop" {
				listen.Close()
				// os.Exit(1)
			}
		}
	}()

	for {
		conn, _ := listen.AcceptTCP()
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		clients = append(clients, conn)

		go handleConnection(conn)
	}
}
