package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run benchmark_server port")
		os.Exit(1)
	}

	port := os.Args[1]
	port = ":" + port

	//listen on socket ln
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}

	defer ln.Close()
	//accept connection and create a new socket
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		defer conn.Close()

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			os.Exit(1)
		}
		conn.Write([]byte(message + "\n"))
		conn.Close()
	}

}
