package main

import (
	"fmt"

	"example.com/m/server"
)

func main() {
	fmt.Println("Start")
	ser := server.NewServer("127.0.0.1", 8888)
	ser.Start()
}
