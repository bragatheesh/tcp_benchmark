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
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run benchmark_client host port num_tests")
		os.Exit(1)
	}

	var totalElapsed time.Duration
	var tempSize int
	var PktIOIntf PktIO
	var mySock PktIOSocket
	var tmpPort int
	var err error
	var sock uint32
	buffer := []byte("Hello from the client\n")
	totalSentSize := 0
	totalRecvdSize := 0
	loops := 0

	numTests, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Error converting "+os.Args[3]+" to int", err.Error())
		os.Exit(1)
	}

	mySock.Type = SOCKET_TCP
	tmpPort, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Unable to convert "+os.Args[2]+" to int ", err.Error())
		os.Exit(1)
	}

	mySock.RemoteL4Port = uint16(tmpPort)
	mySock.RemoteIP = os.Args[1]

	sock, err = PktIOIntf.CreateSocket(mySock)
	if err != nil {
		fmt.Println("Unable to create socket ", err.Error())
		os.Exit(1)
	}

	defer PktIOIntf.DeleteSocket(sock)

	for loops < numTests {
		var recvBuffer = make([]byte, 1024)

		err = PktIOIntf.Connect(sock, mySock, time.Duration(0))
		if err != nil {
			fmt.Println("Error connecting", err.Error())
			os.Exit(1)
		}

		start := time.Now()
		tempSize, err = PktIOIntf.Write(buffer)
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}

		_, err = PktIOIntf.Read(recvBuffer)
		elapsed := time.Since(start)
		totalElapsed += elapsed
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			os.Exit(1)
		}

		//just to test the APIS
		n := bytes.IndexByte(recvBuffer, 0)
		s := string(recvBuffer[:n])
		fmt.Println("buffer from server ", s)
		totalSentSize += tempSize
		totalRecvdSize += n
		PktIOIntf.DeleteSocket(sock)
		loops++
	}

	fmt.Printf("PASS\nResults:\nBytes sent: %d\nBytes recvd: %d\nTotal time:"+
		"%s\nAvgRTT: %f Seconds %d Nanoseconds\n", totalSentSize,
		totalRecvdSize, totalElapsed,
		totalElapsed.Seconds()/float64(numTests),
		totalElapsed.Nanoseconds()/int64(numTests))

}
