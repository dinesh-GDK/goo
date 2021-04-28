package main

import (
	"fmt"
)

func command_palette(command string) {

	switch command {
	case ":clear":
		fmt.Println("\033[2J\033[H")

	case ":myip":
		ip_addr, err := get_ip()
		if err != nil {
			return
		}
		fmt.Println("-->> YOUR IP - " + ip_addr + " <<--")

	default:
		fmt.Println("-->> WRONG COMMAND <<--")
	}
}
