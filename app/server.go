package main

import (
	"fmt"
	"net"
	"os"
)

const PORT = "4221"

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+PORT)
	if err != nil {
		fmt.Printf("Failed to bind to port %s: %v\n", PORT, err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Listening on port %s\n", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	requestBuffer := make([]byte, 1024)

	_, err := conn.Read(requestBuffer)
	if err != nil {
		fmt.Printf("Error reading connection: %v\n", err)
		return
	}

	req := ParseRequest(string(requestBuffer))
	resp := processRequest(req)
	conn.Write([]byte(resp))
}

func processRequest(req Request) string {
	handler := GetHandler(req.Endpoint)
	if handler == nil {
		return NotFoundResponse()
	}

	return handler(req)
}
