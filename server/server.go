package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// 模块初始化函数 import 包时被调用
func init() {
	fmt.Println("Server module init function")
}

type Server struct {
	Ip string
	Port int

	// onlineMap
	OnlineMap  map[string]*User
	mapLock sync.RWMutex

	// broadcast channel
	Message chan string
	UserStatus chan string
}

func NewServer(ip string, port int) (*Server) {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
		UserStatus: make(chan string),
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
	}
}

func(this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func(this *Server) AddUser(user *User) {
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock();
	this.Broadcast(user, " ON line!!!")
}

func(this *Server) DelUser(user *User) {
	this.mapLock.Lock()
	delete(this.OnlineMap, user.Name)
	this.mapLock.Unlock()
	this.Broadcast(user, " OFF line!!!")
	user.Offline()
}

func (this *Server) Handler(conn net.Conn) {
	// current connect work
	fmt.Println("connect success")

	user := NewUser(conn, this)

	fmt.Println("user name is ", user.Name, " connect success")

	// user online
	this.AddUser(user)

	isLive := make(chan bool)

	// receive user's send msg
	go func ()  {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n==0 {
				user.Offline()
				//这样的写法会造成服务端出现大量CLOSE_WAIT，return也只是退出当前协程，而不是Handle
				return
			}

			//读到末尾的时候err是io.EOF
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			
			//去除最后的'\n'并转为字符串
			msg := string(buf[:n-1])
			user.OnMessage(msg)

			isLive <- true
		}
	}()

	// 空select会一直阻塞
	for {
		select {
		case <-isLive:
		//当前用户是活跃的,不做任何事情，select中会执行下面这句重置
    //time.After(time.Second * 10)重新执行即刻重置定时器，定时到后会发送信息
    //这里select 第二个case定时器触发后，处于阻塞状态。当满足第一个 case 的条件后，
    //打破了 select 的阻塞状态，每个条件又开始判断，第2个 case 的判断条件一执行，就重置定时器了。
		case <-time.After(time.Second*10):
			this.DelUser(user)

			return // runtime.Goexit()
		}
	}

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