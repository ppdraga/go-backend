package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var nickname *string

func main() {
	nickname = flag.String("nickname", "Unknown", "Nickname for chatting")
	flag.Parse()
	log.Print("Nickname: ", *nickname)

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send auth info to server
	authInfo := "Set nickname: " + *nickname + "\n"
	_, err = io.WriteString(conn, authInfo)
	if err != nil {
		log.Print(err)
	}

	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
