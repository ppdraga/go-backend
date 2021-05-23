package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	clientConnections = make(map[net.Conn]bool)
	sendMsgMutex      = make(chan struct{}, 1)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	messages := make(chan string)
	go consoleReader(messages)
	go broadcaster(messages)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		clientConnections[conn] = true
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		sendMsgMutex <- struct{}{}
		msg := time.Now().Format("15:04:05\n\r")
		_, err := io.WriteString(c, msg)
		<-sendMsgMutex
		if err != nil {
			delete(clientConnections, c)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func broadcaster(messages <-chan string) {
	for {
		msg := <-messages
		log.Printf("Got message from console: %s", msg)
		for clientConn := range clientConnections {
			sendMsgMutex <- struct{}{}
			log.Printf("Sending message to client: %s", clientConn.RemoteAddr().String())
			_, err := io.WriteString(clientConn, msg+"\n")
			<-sendMsgMutex
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func consoleReader(messages chan<- string) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		messages <- input.Text()
	}
}
