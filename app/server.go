package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	requestBuffer := make([]byte, 1024)
	response := ""
	_, err := conn.Read(requestBuffer)
	if err != nil {
		fmt.Println("Error in reading connection: ", err.Error())
		os.Exit(1)
	}

	requestParts := strings.Split(string(requestBuffer), "\r\n")
	endpoint := strings.Split(requestParts[0], " ")[1]

	if endpoint == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(endpoint, "/echo/") {
		echo := strings.Split(endpoint, "/")[2]
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echo), echo)
	} else if strings.HasPrefix(endpoint, "/user-agent") {
		userAgentHeader := ""
		for _, header := range requestParts {
			if strings.HasPrefix(strings.ToLower(header), "user-agent") {
				userAgentHeader = header
			}
		}
		userAgent := strings.Split(userAgentHeader, " ")[1]
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	conn.Write([]byte(response))
}
