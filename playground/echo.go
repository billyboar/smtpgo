package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Program started")
	ln, err := net.Listen("tcp", ":25")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	log.Printf("Connection from %v is established.", c.RemoteAddr())
	buf := make([]byte, 4096)
	_, err := c.Write([]byte("220 hello man \n"))
	if err != nil {
		c.Close()
	}
	for {
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			c.Close()
			break
		}
		checker(string(buf[0:n-2]), c)
		if err != nil {
			c.Close()
			break
		}
	}
	log.Printf("Connection from %v is closed.", c.RemoteAddr())
}

func checker(command string, c net.Conn) {
	//fmt.Println(command)
	if command == "EHLO" {
		n, err := c.Write([]byte("Hello Yourself\n"))
		if err != nil || n == 0 {
			c.Close()
		}
	}
}
