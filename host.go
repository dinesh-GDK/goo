package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var users_list = &User{}
var trash = &User{name: "trash"}

func broadCast(message string, exclude_user *User) {

	fmt.Print(message)

	dummy := users_list.next
	for {

		if(dummy == nil) {
			break
		}

		if(dummy == exclude_user) {
			dummy = dummy.next
			continue
		}

		_, err := dummy.conn.Write([]byte(message))
		error_handler(err)

		dummy = dummy.next
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

		exist := users_list.search_user_name(name)

		if !exist {
			_, err = conn.Write([]byte(string("set\n")))
			error_handler(err)			
			break

		} else {
			_, err = conn.Write([]byte(string("no_set\n")))
			error_handler(err)
		}
	}

	user := &User{
		name: name,
		conn: conn,
	}

	users_list.insert(user)

	clear_chat_line(users_list.name)
	broadCast("-->> "+name+" joined <<--\n", trash)
	print_chat_line(users_list.name)

	// read from client
	for {
		client_msg, err := bufio.NewReader(conn).ReadString('\n')
		// when connection with client is broken
		if err != nil {
			break
		}

		clear_chat_line(users_list.name)
		broadCast(name+" >> "+client_msg, user)
		print_chat_line(users_list.name)
	}
	
	users_list.remove(user)

	users_list.print()

	clear_chat_line(users_list.name)
	broadCast("-->> "+name+" left <<--\n", trash)
	print_chat_line(users_list.name)
	
	conn.Close()
	users_list.remove(user)
	users_list.print()
}

func host() {

	var PORT string

	// check for port validity
	// fmt.Print("Enter the PORT to host: ")
	// fmt.Scanf("%s", &PORT)
	PORT = "5000"

	addr, err := net.ResolveTCPAddr("tcp4", ":"+PORT)
	error_handler(err)

	listen, err := net.ListenTCP("tcp4", addr)
	error_handler(err)

	defer listen.Close()

	fmt.Println("Hosting at: " + get_ip() + ":" + PORT)

	// fmt.Print("Enter User name: ")
	// fmt.Scanf("%s", &users_list.name)
	users_list.name = "qwe"

	// send message
	go func() {
		for {

			print_chat_line(users_list.name)

			text, err := bufio.NewReader(os.Stdin).ReadString('\n')
			error_handler(err)
			temp := strings.TrimSpace(string(text))

			if temp == ":stop" {
				listen.Close()
			}

			clear_chat_prev_line()
			broadCast(users_list.name+" >> "+text, trash)
		}
	}()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			clear_chat_line(users_list.name)
			fmt.Println("-->> DISCONNECTED <<--")
			os.Exit(34)
		}

		go handleConnection(conn)
	}
}
