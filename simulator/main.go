package main

import (
	"fmt"
	"net"
	"time"

	"./data"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3500")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tcpConn := conn.(*net.TCPConn)
	err = tcpConn.SetWriteBuffer(1024)
	if err != nil {
		fmt.Println(err.Error())
	}
	tcpConn.Close()
}
