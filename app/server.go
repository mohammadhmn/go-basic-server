package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type RequestHandler func(requestParts []string) string

const PORT = "4221"

const (
	OK          = "HTTP/1.1 200 OK"
	BAD_REQUEST = "HTTP/1.1 400 Bad Request"
	NOT_FOUND   = "HTTP/1.1 404 Not Found"
)

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

	response := processRequest(string(requestBuffer))
	conn.Write([]byte(response))
}

func processRequest(request string) string {
	requestParts := strings.Split(request, "\r\n")
	if len(requestParts) == 0 {
		return badRequestResponse()
	}

	endpoint := strings.Split(requestParts[0], " ")[1]

	handler := getHandler(endpoint)
	if handler == nil {
		return notFoundResponse()
	}

	return handler(requestParts)
}

func getHandler(endpoint string) RequestHandler {
	switch {
	case endpoint == "/":
		return HandleRoot
	case strings.HasPrefix(endpoint, "/echo/"):
		return HandleEcho
	case strings.HasPrefix(endpoint, "/user-agent"):
		return HandleUserAgent
	case strings.HasPrefix(endpoint, "/files"):
		return HandleFile
	default:
		return nil
	}
}

func HandleRoot(requestParts []string) string {
	return okResponse()
}

func HandleEcho(requestParts []string) string {
	endpoint := strings.Split(requestParts[0], " ")[1]
	echo := strings.TrimPrefix(endpoint, "/echo/")
	return responseBuilder(OK, "text/plain", len(echo), echo)
}

func HandleUserAgent(requestParts []string) string {
	for _, header := range requestParts {
		if strings.HasPrefix(strings.ToLower(header), "user-agent") {
			userAgent := strings.Join(strings.Split(header, " ")[1:], " ")
			return responseBuilder(OK, "text/plain", len(userAgent), userAgent)
		}
	}
	return badRequestResponse()
}

func HandleFile(requestParts []string) string {
	directory := os.Args[2]
	if directory == "" {
		fmt.Println("Error reading file directory")
		return badRequestResponse()
	}

	endpoint := strings.Split(requestParts[0], " ")[1]
	file := strings.TrimPrefix(endpoint, "/files/")
	fileData, err := os.ReadFile(filepath.Join(directory, file))
	if err != nil {
		return notFoundResponse()
	}
	return responseBuilder(OK, "application/octet-stream", len(fileData), string(fileData))
}

func responseBuilder(statusLine string, contentType string, contentLength int, body string) string {
	return fmt.Sprintf("%s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", statusLine, contentType, contentLength, body)
}

func okResponse() string {
	return OK + "\r\n\r\n"
}

func badRequestResponse() string {
	return BAD_REQUEST + "\r\n\r\n"
}

func notFoundResponse() string {
	return NOT_FOUND + "\r\n\r\n"
}
