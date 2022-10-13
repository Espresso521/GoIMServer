package main

import (
	"fmt"

	"example.com/myfistgo/m/instant-chat/src/server"
)

// 两个包都属于main包，没必要import
func main() {
	fmt.Println("Start")
	ser := server.NewServer("127.0.0.1", 8888)
	ser.Start()
}

