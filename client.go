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
	if err != nil {
		fmt.Println("1")
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("DISCONNECTED")
	}

	var user_name string
	for {

		fmt.Print("Enter User name: ")
		fmt.Scanf("%s", &user_name)

		_, err = conn.Write([]byte(user_name + "\n"))
		if err != nil {
			fmt.Println("2")
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("3")
		}
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
			if err != nil {
				fmt.Println("4")
			}

			fmt.Print("\033[2K")
			fmt.Printf("\033[%dD", len(user_name)+6)
			fmt.Print(message)
			fmt.Printf("%s -->> ", user_name)
		}
	}()

	// send message
	for {

		fmt.Printf("%s -->> ", user_name)

		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("5")
		}
		temp := strings.TrimSpace(string(text))

		fmt.Print("\033[A")
		fmt.Print("\033[K")

		fmt.Print(user_name + " >> " + text)

		if temp == ":stop" {
			conn.Close()
			// time.Sleep(2 * time.Second)
			os.Exit(3)
		}

		_, err = conn.Write([]byte(string(text)))
		if err != nil {
			fmt.Println("6")
		}

	}

}
