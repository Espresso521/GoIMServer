package server

import (
	"fmt"
	"net"
	"sync"

	"example.com/m/users"
)

// 模块初始化函数 import 包时被调用
func init() {
	fmt.Println("Server module init function")
}

type Server struct {
	Ip string
	Port int

	// onlineMap
	OnlineMap  map[string]*users.User
	mapLock sync.RWMutex

	// broadcast channel
	Message chan string
}

func NewServer(ip string, port int) (*Server) {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*users.User),
		Message:make(chan string),
	}

	fmt.Printf("ip is %s, port is %d", ip, port)
	return server
}

func (this *Server) ListenMsg() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock();

		fmt.Println("send success")
	}
}

func(this *Server) Broadcast(user *users.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// current connect work
	fmt.Println("connect success")

	user := users.NewUser(conn)

	fmt.Println("user name is ", user.Name)
	// user online
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock();

	// broadcast user online event
	this.Broadcast(user, "is online")

	fmt.Println("Broadcast success")
	// temp block
	select {}
}

func (this *Server) Start() {
	//socket listen
	address := fmt.Sprint(this.Ip,":",this.Port)
	fmt.Println("server address is :", address)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	//close listen socket
	defer listener.Close()

	//start listener msg
	go this.ListenMsg()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Listen err:", err)
			continue
		}

		go this.Handler(conn)
	}

	//accept

	//do handler

	
}