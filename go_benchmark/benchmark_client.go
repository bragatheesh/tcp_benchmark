package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run benchmark_client host port num_tests")
		os.Exit(1)
	}

	var totalElapsed time.Duration
	var tempSize int
	totalSentSize := 0
	totalRecvdSize := 0
	loops := 0

	numTests, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Error converting "+os.Args[3]+" to int", err.Error())
		os.Exit(1)
	}
	port := os.Args[2]
	host := os.Args[1]
	connection := host + ":" + port

	for loops < numTests {
		conn, err := net.Dial("tcp", connection)
		if err != nil {
			fmt.Println("Error connecting", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		text := "Hello from the client"

		start := time.Now()
		tempSize, err = fmt.Fprint(conn, text+"\n")

		message, err := bufio.NewReader(conn).ReadString('\n')
		elapsed := time.Since(start)
		totalElapsed += elapsed
		if err != nil {
			fmt.Println("Error", err.Error())
			os.Exit(1)
		}
		totalSentSize += tempSize
		totalRecvdSize += len(message)

		//fmt.Println("Message from server: " + message)
		conn.Close()
		loops++
	}

	fmt.Printf("Results:\nBytes sent: %d\nBytes recvd: %d\nTotal time: %s\nAvgRTT: %f Seconds %d Nanoseconds\n", totalSentSize, totalRecvdSize, totalElapsed,
		totalElapsed.Seconds()/float64(numTests), totalElapsed.Nanoseconds()/int64(numTests))

}
