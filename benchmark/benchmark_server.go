package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	SOCKET_RAW uint8 = iota
	SOCKET_IP
	SOCKET_TCP
	SOCKET_UDP
)

type Endpoint struct {
	Vrf      string
	IPAdress string
	L4Port   uint16
	Zone     string
}

type PktIOSocket struct {
	Type         uint8
	Vrf          uint16
	Port         uint16
	Vlan         uint16
	EthType      uint16
	IPAddress    string
	RemoteIP     string
	L3Proto      uint16
	L4Port       uint16
	RemoteL4Port uint16
}

type PktIO interface {
	CreateSocket(PktIOSocket) (uint32, error)
	DeleteSocket(uint32) error
	SetSockoptTCPMD5(uint32, string, string)
	Connect(uint32, PktIOSocket, time.Duration) error
	Listen(uint32) error
	Accept(uint32) (uint32, error)
	Read([]byte) (int, error)
	ReadFrom([]byte) (int, Endpoint, error)
	Write([]byte) (int, error)
	WriteTo([]byte, Endpoint) (int, error)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run benchmark_server port")
		os.Exit(1)
	}

	var PktIOIntf PktIO
	var mySock PktIOSocket
	var sock uint32
	var err error
	var tmpPort int

	mySock.Type = SOCKET_TCP
	tmpPort, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Unable to convert "+os.Args[1]+" to int ", err.Error())
		os.Exit(1)
	}

	mySock.L4Port = uint16(tmpPort)

	sock, err = PktIOIntf.CreateSocket(mySock)
	if err != nil {
		fmt.Println("Unable to create socket ", err.Error())
		os.Exit(1)
	}

	defer PktIOIntf.DeleteSocket(sock)

	err = PktIOIntf.Listen(sock)
	if err != nil {
		fmt.Println("Unable to put socket in listen mode ", err.Error())
		os.Exit(1)
	}
	//accept connection and create a new socket
	for {

		buffer := make([]byte, 1024)

		conn, err := PktIOIntf.Accept(sock)
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		defer PktIOIntf.DeleteSocket(conn)

		//just about to implement read and write

		_, err = PktIOIntf.Read(buffer)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			os.Exit(1)
		}

		_, err = PktIOIntf.Write(buffer)
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}

		//just to test the APIS
		n := bytes.IndexByte(buffer, 0)
		s := string(buffer[:n])
		fmt.Println("buffer from client ", s)

		PktIOIntf.DeleteSocket(conn)
	}

}
