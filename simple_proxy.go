package main

import (
	"io"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Proxy started on :8080,forwarding to google.com:80")

	for {
		client.Conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConnection(clientConn)
	}

}

func handleConnection(clientConn net.Conn) {

	targetConn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		log.Println("Dial error:", err)
		clientConn.Close()
		return
	}

	go func() {
		io.Copy(targetConn, clientConn)
		targetConn.Close()
	}()

	io.Copy(clientConn, targetConn)
	clientConn.Close()

}
