package main

import (
	"fmt"
	"os"
)

func main() {

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

		if activity == "h" {
			host()
			break

		} else if activity == "j" {
			client()
			break

		} else {
			fmt.Println("Invalid input")
		}
	}

}
