package main

import (
	"fmt"
	"net"
)

type User struct {
	name string
	conn *net.TCPConn
	next *User
}

func (head *User) remove(user *User) {
	
	dummy := head
	for {
		if(dummy.next == nil) {
			break
		}

		if(dummy.next == user) {
			dummy.next = dummy.next.next
			break
		}

		dummy = dummy.next
	}
}

func (head *User) insert(user *User) {

	dummy := head
	for {
		if(dummy.next == nil) {
			break
		}

		dummy = dummy.next
	}

	dummy.next = user
}

func (head *User) search(name string) bool {

	dummy := head
	for {
		if(dummy == nil) {
			break
		}

		if(dummy.name == name) {
			return true
		}

		dummy = dummy.next
	}

	return false
}

func (head *User) print() {

	dummy := head
	for {

		fmt.Println(dummy.name)

		if(dummy.next == nil) {
			break
		}
		
		dummy = dummy.next
	}
}
