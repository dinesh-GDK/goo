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

func (head *User) search_user_name(name string) bool {

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

// func main() {
// 	head := &User{name: "head"}

// 	// dummy := head

// 	// dummy.name = "a"
// 	// dummy.next = &User{}
// 	// dummy = dummy.next

// 	// dummy.name = "b"
// 	// dummy.next = &User{}
// 	// dummy = dummy.next

// 	// dummy.name = "c"
// 	// dummy.next = &User{}
// 	// dummy = dummy.next

// 	// head.print()
// 	// fmt.Println()

// 	// head.delete(head.next)
// 	// head.print()

// 	head.print()

// 	head.insert(&User{name: "a"})
// 	head.insert(&User{name: "b"})
// 	head.insert(&User{name: "c"})
// 	head.print()

// 	head.delete(head.next)
// 	head.print()

	
	

// }