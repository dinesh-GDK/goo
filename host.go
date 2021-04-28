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
var CIPHER string
var CIPHER_ATTEMPTS int

func host() {

	var PORT string

	if !DEBUG {
		// check for port validity
		fmt.Print("Enter the PORT to host: ")
		fmt.Scanf("%s", &PORT)

	} else {
		PORT = "5000"
	}

	addr, err := net.ResolveTCPAddr("tcp4", ":"+PORT)
	error_handler(err)

	listen, err := net.ListenTCP("tcp4", addr)
	error_handler(err)

	defer listen.Close()

	ip_addr, err := get_ip()
	if err != nil {
		return
	}
	fmt.Println("Hosting at: " + ip_addr + ":" + PORT)

	var host_user_name string
	if !DEBUG {
		fmt.Print("Enter User name: ")
		fmt.Scanf("%s", &host_user_name)

		fmt.Print("Set CIPHER: ")
		fmt.Scanf("%s", &CIPHER)

		fmt.Print("Set number of CIPHER attempts: ")
		fmt.Scanf("%d", &CIPHER_ATTEMPTS)

	} else {
		host_user_name = "qwe"
		CIPHER = "aaa"
		CIPHER_ATTEMPTS = 1
	}

	users_list.init(host_user_name)
	ban_users_list.init("atrash")

	end := make(chan bool)

	go host_send_message(listen, end)
	go host_handle_multiple_clients(listen)

	<-end
}

func host_send_message(listen *net.TCPListener, end chan bool) {

	for {
		print_chat_line(users_list.head.name)

		host_message, err := bufio.NewReader(os.Stdin).ReadString('\n')
		error_handler(err)
		temp := strings.TrimSpace(string(host_message))

		if temp == ":stop" {
			listen.Close()
			fmt.Print("--> STOPPED HOSTING <--\n")
			end <- true
			return
		}

		if temp[:3] == ":ub" {

			unban_num := -1
			if len(temp) > 4 {
				unban_num, err = strconv.Atoi(temp[4:])

				if err != nil || unban_num > ban_users_list.length {
					fmt.Println("--> WRONG BAN NUMBER <--")
					continue
				}
			}

			dummy := ban_users_list.head.next
			if dummy == nil {
				fmt.Println("--> No User(s) banned <--")
			}

			i := 1
			for {
				if dummy == nil {
					break
				}

				if unban_num < 0 {
					fmt.Println(strconv.Itoa(i) + ") " + dummy.conn.RemoteAddr().(*net.TCPAddr).IP.String())
				}

				if unban_num == i {
					ban_users_list.remove(dummy)
				}

				dummy = dummy.next
			}

			continue
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

	for i := 1; i <= CIPHER_ATTEMPTS; i++ {

		cipher, err := bufio.NewReader(conn).ReadString('\n')
		error_handler(err)

		cipher = strings.TrimSuffix(cipher, "\n")

		if cipher == CIPHER {
			_, err := conn.Write([]byte(string(":m\n")))
			error_handler(err)
			return true

		} else if i == CIPHER_ATTEMPTS {
			_, err := conn.Write([]byte(string(":le\n")))
			error_handler(err)
			break

		} else {
			_, err := conn.Write([]byte(string(":nm#" + strconv.Itoa(CIPHER_ATTEMPTS-i) + "\n")))
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
		fmt.Print("*")

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
