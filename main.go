package main

import (
	"fmt"

	"example.com/m/server"
)

func main() {
	fmt.Println("Start")
	ser := server.NewServer("0.0.0.0", 5210)
	ser.Start()
}
