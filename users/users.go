package users

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
}

// create a user
func NewUser(conn net.Conn) *User {
	userAddr:=conn.RemoteAddr().String()

	user:=&User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
	}

	return user
}

// listen channel
func (this *User) ListenMsg() {
	for{
		msg := <-this.C
		fmt.Println("msg success is ", msg)
		this.conn.Write([]byte(msg+"\n"))
	}
}