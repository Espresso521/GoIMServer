package server

import (
	"net"
	"strings"
	"sync"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *IMServer
	once sync.Once
}

// create a user
func NewUser(conn net.Conn, server *IMServer) *User {

	userName:= strings.Split(conn.RemoteAddr().String(),":")[1] 

	user:=&User{
		Name: userName,
		Addr: conn.RemoteAddr().String(),
		C: make(chan string),
		conn: conn,
		server: server,
	}

	// goroutine for user to listen msg
	go user.ListenMsg()

	return user
}

func (this *User) SafeClose() {
	this.once.Do(func() {
		close(this.C)
		this.conn.Close()
	})
}

func (this *User) Offline() {
		// refer to https://go101.org/article/channel-closing.html
		this.SafeClose()
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) OnMessage(msg string) {
	this.server.DispathMsg(this, msg)
}

// listen channel
func (this *User) ListenMsg() {
  //当u.C通道关闭后，不再进行监听并写入信息
	for msg := range this.C {
			_, err := this.conn.Write([]byte(msg + "\n"))
			if err != nil {
				panic(err)
			}
	}
}