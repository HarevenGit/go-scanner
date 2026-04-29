package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	timeout := 2 * time.Second
	Ip := "192.138.0.1:5050"
	conn, err := net.DialTimeout("tcp", Ip, timeout)
	if err != nil {
		fmt.Println(err, "err")
	} else {
		fmt.Println(conn, "sucess")
		conn.Close()
	}
}
