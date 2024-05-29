package main

import (
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) {
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
