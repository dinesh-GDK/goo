package main

import (
	"fmt"
)

func command_palette(command string) {

	switch command {
	case ":clear":
		fmt.Println("\033[2J\033[H")

	case ":myip":
		fmt.Println("-->> YOUR IP - " + get_ip() + " <<--")

	default:
		fmt.Println("-->> WRONG COMMAND <<--")
	}
}
