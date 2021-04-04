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

func (list *Users_list) init(head_name string) {
	list.head = &User{name: head_name}
	list.tail = list.head
}

func (list *Users_list) insert(user *User) {

	list.tail.next = user
	list.tail = list.tail.next
	list.length++
}

func (list *Users_list) remove(user *User) {
	
	dummy := list.head
	for {
		if dummy.next == nil {
			break
		}

		if dummy.next == user {
			dummy.next = dummy.next.next
			break
		}

		dummy = dummy.next
	}
}

func (list *Users_list) search_by_user_name(name string) bool {

	dummy := list.head
	for {
		if dummy == nil {
			break
		}

		if dummy.name == name {
			return true
		}

		dummy = dummy.next
	}

	return false
}

func (list *Users_list) search_by_conn(conn *net.TCPConn) bool {

	dummy := list.head.next
	for {
		if dummy == nil {
			break
		}

		if dummy.conn.RemoteAddr().(*net.TCPAddr).IP.String() == conn.RemoteAddr().(*net.TCPAddr).IP.String() {
			return true
		}

		dummy = dummy.next
	}

	return false
}

func (list *Users_list) print() {

	dummy := list.head
	for {

		fmt.Println(dummy)

		if dummy.next == nil {
			break
		}
		
		dummy = dummy.next
	}
}
