package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// first user is the host
var users_list = &Users_list{}
var ban_users_list = &Users_list{}
var trash = &User{name: "trash"}

var ROOM_CAPACITY int = 10
var CIPHER string = "aaa"
var CIPHER_ATTEMPTS int = 2

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
	users_list.init("qwe")
	ban_users_list.init("atrash")

	// host_set_cipher()

	go host_send_message(listen)
	host_handle_multiple_clients(listen)
}

func host_set_cipher() {
	fmt.Print("Set CIPHER: ")
	fmt.Scanf("%s", &CIPHER)

	fmt.Print("Set number of CIPHER attempts: ")
	fmt.Scanf("%i", &CIPHER_ATTEMPTS)
}

func host_send_message(listen *net.TCPListener) {

	for {
		print_chat_line(users_list.head.name)

		host_message, err := bufio.NewReader(os.Stdin).ReadString('\n')
		error_handler(err)
		temp := strings.TrimSpace(string(host_message))

		if temp == ":stop" {
			listen.Close()
		}

		if temp[0] == ':' {
			command_palette(temp)
			continue
		}

		host_broadcast(users_list.head.name+" >> "+host_message, trash, true)
	}
}

func host_handle_multiple_clients(listen *net.TCPListener) {

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			clear_chat_line(users_list.head.name)
			fmt.Println("-->> DISCONNECTED <<--")
			return
		}

		if ban_users_list.search_by_conn(conn) {
			_, err := conn.Write([]byte(":b\n"))
			error_handler(err)

		} else {
			_, err := conn.Write([]byte(":nb\n"))
			error_handler(err)
			go host_handle_client(conn)
		}
	}
}

func host_handle_client(conn *net.TCPConn) {

	if users_list.length >= ROOM_CAPACITY {
		conn.Write([]byte(string(":f\n")))
		return
	}

	if !host_cipher_verification(conn) {
		return
	}

	user := &User{
		name: host_set_client_user_name(conn),
		conn: conn,
	}

	users_list.insert(user)
	host_broadcast("-->> "+user.name+" joined <<--\n", trash, false)

	// read from client
	for {
		client_message, err := bufio.NewReader(conn).ReadString('\n')
		// when connection with client is broken
		if err != nil {
			break
		}

		host_broadcast(user.name+" >> "+client_message, user, false)
	}

	host_broadcast("-->> "+user.name+" left <<--\n", trash, false)

	conn.Close()
	users_list.remove(user)
}

func host_cipher_verification(conn *net.TCPConn) bool {

	for i := 0; i < CIPHER_ATTEMPTS; i++ {

		cipher, err := bufio.NewReader(conn).ReadString('\n')
		error_handler(err)

		cipher = strings.TrimSuffix(cipher, "\n")

		if cipher == CIPHER {
			_, err := conn.Write([]byte(string(":m\n")))
			error_handler(err)
			return true

		} else if i == CIPHER_ATTEMPTS-1 {
			_, err := conn.Write([]byte(string(":lie\n")))
			error_handler(err)
			break

		} else {
			_, err := conn.Write([]byte(string(":nm#" + strconv.Itoa(CIPHER_ATTEMPTS-i-1) + "\n")))
			error_handler(err)
		}
	}

	ban_users_list.insert(&User{conn: conn})
	return false
}

func host_set_client_user_name(conn *net.TCPConn) string {

	var name string
	var err error

	// set user name
	for {

		name, err = bufio.NewReader(conn).ReadString('\n')
		error_handler(err)

		name = strings.TrimSuffix(name, "\n")

		exist := users_list.search_by_user_name(name)

		if !exist {
			_, err = conn.Write([]byte(string(":s\n")))
			error_handler(err)
			break

		} else {
			_, err = conn.Write([]byte(string(":ns\n")))
			error_handler(err)
		}
	}

	return name
}

func host_broadcast(message string, exclude_user *User, self bool) {

	if self {
		clear_chat_prev_line()
	} else {
		clear_chat_line(users_list.head.name)
	}

	fmt.Print(message)

	dummy := users_list.head.next
	for {

		if dummy == nil {
			break
		}

		if dummy == exclude_user {
			dummy = dummy.next
			continue
		}

		_, err := dummy.conn.Write([]byte(message))
		error_handler(err)

		dummy = dummy.next
	}

	if !self {
		print_chat_line(users_list.head.name)
	}
}
