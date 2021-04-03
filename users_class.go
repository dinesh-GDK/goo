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

type Users_list struct {
	head *User
	tail *User
	length int
}

func (list *Users_list) remove(user *User) {
	
	dummy := list.head
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

func (list *Users_list) insert(user *User) {

	list.tail.next = user
	list.tail = list.tail.next
	list.length++
}

func (list *Users_list) search(name string) bool {

	dummy := list.head
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

func (list *Users_list) print() {

	dummy := list.head
	for {

		fmt.Println(dummy.name)

		if(dummy.next == nil) {
			break
		}
		
		dummy = dummy.next
	}
}

// func main() {

// 	l := &Users_list{}
// 	l.head = &User{name: "asd"}

// 	l.tail = l.head

// 	a := &User{name: "a"}
// 	b := &User{name: "b"}
// 	c := &User{name: "c"}
// 	d := &User{name: "d"}
// 	e := &User{name: "e"}

// 	l.insert(a)
// 	l.insert(b)
// 	l.insert(c)
// 	l.insert(d)
// 	l.insert(e)
// 	fmt.Println(l.length)
// 	l.print()

// 	l.remove(c)
// 	l.print()

// 	fmt.Println(l.search("d"))
// 	fmt.Println(l.search("z"))
// }
