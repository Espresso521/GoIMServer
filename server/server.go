package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

// 模块初始化函数 import 包时被调用
func init() {
	fmt.Println("Server module init function")
}

type IMServer struct {
	Ip string
	Port int

	// onlineMap
	OnlineMap  map[string]*User
	mapLock sync.RWMutex

	// broadcast channel
	Message chan []byte
	UserStatus chan string

	// socket.io websocket server
	Server *socketio.Server 
}

// 定义结构体
type DispatchMsg struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewServer(ip string, port int) (*IMServer) {
	server := &IMServer{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan []byte),
		UserStatus: make(chan string),
		Server : socketio.NewServer(nil),
	}

	return server
}

func (this *IMServer) ListenMsg() {
	for {
		msg := <-this.Message
		dmsg := DispatchMsg {}
		err := json.Unmarshal(msg, &dmsg)

		if err != nil {
			fmt.Println("err is " + err.Error())
			continue
		}

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			if cli.Name != dmsg.Name {
				cli.C <- dmsg.Content
			}
		}
		this.mapLock.Unlock();
	}
}

func (this *IMServer) ListenStatus() {
	for {
		msg := <-this.UserStatus

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock();
	}
}

func(this *IMServer) BroadcastStatus(user *User, msg string) {
	sendMsg := "[" + user.Name  + "]" + ":" + msg
	this.UserStatus <- sendMsg
}

func(this *IMServer) DispathMsg(user *User, msg string) {
	sendMsg := "[" + user.Name + "]" + ":" + msg
  // 序列化
	buf, err := json.Marshal(DispatchMsg{
		Name:user.Name,
		Content: sendMsg,
	})

	if err != nil {
		fmt.Println("err is " + err.Error())
		return
	}

	this.Message <- buf
}

func(this *IMServer) AddUser(user *User) {
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock();
	this.BroadcastStatus(user, " ON line!!!")
}

func(this *IMServer) DelUser(user *User) {
	this.mapLock.Lock()
	delete(this.OnlineMap, user.Name)
	this.mapLock.Unlock()
	this.BroadcastStatus(user, " OFF line!!!")
	user.Offline()
}

func (this *IMServer) Handler(conn socketio.Conn) {
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
				this.DelUser(user)
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
		case <-time.After(time.Second*120):
			this.DelUser(user)

			return // runtime.Goexit()
		}
	}

}

func (this *IMServer) Start() {

	//socket listen
	address := fmt.Sprint(this.Ip,":",this.Port)
	fmt.Println("server address is :", address)
	//listener, err := net.Listen("tcp", address)
  router := gin.New()
  log.SetFlags(log.Lshortfile | log.LstdFlags)
	// redis 适配器
	ok, err := this.Server.Adapter(&socketio.RedisAdapterOptions{
			Addr:    "0.0.0.0:5210",
			Prefix:  "kotaku.io",
			Network: "tcp",
	})

	fmt.Println("redis:", ok)

	if err != nil {
			log.Fatal("error:", err)
			return
	}

	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// 连接成功
	this.Server.OnConnect("/", func(conn socketio.Conn) error {
		conn.SetContext("")
		// 申请一个房间
		conn.Join("bcast")
		fmt.Println("连接成功：", conn.ID())
		go this.Handler(conn)
		return nil
	})

	// 接收”bye“事件
	this.Server.OnEvent("/", "bye", func(s socketio.Conn, msg string) string {
		last := s.Context().(string)
		s.Emit("bye", msg)
		fmt.Println("============>", last)
		//s.Close()
		return last
	})
    
	this.Server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
			s.SetContext(msg)
			fmt.Println("=====chat====>", msg)
			return "recv " + msg
	})
	// 连接错误
	this.Server.OnError("/", func(s socketio.Conn, e error) {
			log.Println("连接错误:", e)

	})
	// 关闭连接
	this.Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
			log.Println("关闭连接：", reason)
	})

	go this.Server.Serve()
	defer this.Server.Close()

	//start listener msg
	go this.ListenMsg()

	go this.ListenStatus()

	//router.Use(gin.Recovery(), Cors())
	router.GET("/kotaku.io/*any", gin.WrapH(this.Server))
	router.POST("/kotaku.io/*any", gin.WrapH(this.Server))
	router.StaticFS("/public", http.Dir("../asset"))

	log.Println("Serving at localhost:8000...")
	if err := router.Run(":8000"); err != nil {
			log.Fatal("failed run app: ", err)

	}


	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		fmt.Println("net.Listen err:", err)
	// 		continue
	// 	}

	// 	go this.Handler(conn)
	// }

	//accept

	//do handler

	
}