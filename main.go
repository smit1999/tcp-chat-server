package main

import (
	"log"
	"net"
)

func main() {

	s := servers()

	go s.fetchcommands()

	listener, err := net.Listen("tcp", ":8888")

	if err != nil {
		log.Fatalf("error in connection")
	}

	defer listener.Close()
	log.Printf("started connection on :8888")

	for {
		conn, err := listener.Accept()
		log.Printf("connection details are: %s", conn.LocalAddr())
		if err != nil {
			log.Fatalf("error in establishing connection")
			continue
		}
		go s.newclient(conn)

	}
}
