package main

import (
	"fmt"
	"os"
	"strings"
)

var DEBUG bool = false

func main() {

	if len(os.Args) > 1 && os.Args[1] == "debug" {
		DEBUG = true
	}

	fmt.Println("Hello")

	name, err := os.Hostname()
	error_handler(err)

	fmt.Println("Welcome", name)

	ip_addr := get_ip()
	fmt.Println("Your IP address:", ip_addr)

	for {
		var activity string
		fmt.Print("Do you want to host/join: ")
		fmt.Scanf("%s", &activity)
		activity = strings.ToLower(activity)

		if activity == "h" || activity == "host" {
			host()
			break

		} else if activity == "j" || activity == "join" {
			client()
			break

		} else {
			fmt.Println("Invalid input")
		}
	}

}
