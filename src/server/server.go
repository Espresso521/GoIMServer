package server

import (
	"fmt"
	"net"
)

// 模块初始化函数 import 包时被调用
func init() {
	fmt.Println("Server module init function")
}

type Server struct {
	Ip string
	Port int

	// onlineMap
	OnlineMap map[string]*User


}

func NewServer(ip string, port int) (*Server) {
	server := &Server{
		Ip: ip,
		Port: port,
	}

	fmt.Printf("ip is %s, port is %d", ip, port)
	return server
}

func (this *Server) Handler(conn net.Conn) {
	// current connect work
	fmt.Println("connect success")
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