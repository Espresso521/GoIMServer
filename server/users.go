package server

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

// create a user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr:=conn.RemoteAddr().String()

	user:=&User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}

	// goroutine for user to listen msg
	go user.ListenMsg()

	return user
}

func (this *User) Online() {
	// user online
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock();

	this.server.Broadcast(this, "ON line!")
}

func (this *User) Offline() {
		// user online
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap, this.Name)
		this.server.mapLock.Unlock();
	
		this.server.Broadcast(this, "OFF line!")
}

func (this *User) Kicked() {
	this.SendMsg(this.Name + " are kicked out.")
	close(this.C)
	this.conn.Close()
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) OnMessage(msg string) {
	this.server.Broadcast(this, msg)
}

// listen channel
func (u *User) ListenMsg() {
  //当u.C通道关闭后，不再进行监听并写入信息
	for msg := range u.C {
			_, err := u.conn.Write([]byte(msg + "\n"))
			if err != nil {
				panic(err)
			}
	}
	
	//不监听后关闭conn，conn在这里关闭最合适
	err := u.conn.Close()
	if err != nil {
			panic(err)
	}

}