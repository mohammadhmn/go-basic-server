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
		go HandleConnection(conn)
	}
}
