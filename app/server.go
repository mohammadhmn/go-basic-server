package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type RequestHandler func(request Request) string

const PORT = "4221"

const (
	OK          = "HTTP/1.1 200 OK"
	CREATED     = "HTTP/1.1 201 Created"
	BAD_REQUEST = "HTTP/1.1 400 Bad Request"
	NOT_FOUND   = "HTTP/1.1 404 Not Found"
)

const (
	PlainText   = "text/plain"
	OctetStream = "application/octet-stream"
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

	request := parseRequest(string(requestBuffer))
	response := processRequest(request)
	conn.Write([]byte(response))
}

type Request struct {
	Method   string
	Endpoint string
	Headers  map[string]string
	Body     string
}

func parseRequest(requestString string) Request {
	scanner := bufio.NewScanner(strings.NewReader(requestString))
	scanner.Scan()
	requestLine := scanner.Text()
	requestParts := strings.Split(requestLine, " ")

	method := requestParts[0]
	endpoint := requestParts[1]

	headers := make(map[string]string)
	var body string
	isBody := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isBody = true
			continue
		}
		if isBody {
			body += strings.ReplaceAll(line, "\x00", "")
		} else {
			headerParts := strings.SplitN(line, ": ", 2)
			if len(headerParts) == 2 {
				headers[strings.ToLower(headerParts[0])] = headerParts[1]
			}
		}
	}

	return Request{
		Method:   method,
		Endpoint: endpoint,
		Headers:  headers,
		Body:     body,
	}
}

func processRequest(request Request) string {
	handler := getHandler(request.Endpoint)
	if handler == nil {
		return notFoundResponse()
	}

	return handler(request)
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

func HandleRoot(request Request) string {
	return okResponse()
}

func HandleEcho(request Request) string {
	echo := strings.TrimPrefix(request.Endpoint, "/echo/")
	return responseBuilder(OK, PlainText, len(echo), echo)
}

func HandleUserAgent(request Request) string {
	userAgent, ok := request.Headers["user-agent"]
	if !ok {
		return badRequestResponse()
	}
	return responseBuilder(OK, PlainText, len(userAgent), userAgent)
}

func HandleFile(request Request) string {
	dir, err := readDir()
	if err != nil {
		fmt.Println(err.Error())
		return badRequestResponse()
	}

	filename := strings.TrimPrefix(request.Endpoint, "/files/")

	if request.Method == "GET" {
		fileData, err := os.ReadFile(filepath.Join(dir, filename))
		if err != nil {
			return notFoundResponse()
		}
		return responseBuilder(OK, OctetStream, len(fileData), string(fileData))
	} else if request.Method == "POST" {
		file, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			return notFoundResponse()
		}
		_, err = file.WriteString(request.Body)
		if err != nil {
			return notFoundResponse()
		}
		err = file.Close()
		if err != nil {
			return notFoundResponse()
		}
		return responseBuilder(CREATED, OctetStream, 0, filename)
	}

	return notFoundResponse()
}

func responseBuilder(statusLine string, contentType string, contentLength int, body string) string {
	response := ""
	if statusLine != "" {
		response += statusLine + "\r\n"
	}
	if contentType != "" {
		response += "Content-Type: " + contentType + "\r\n"
	}
	if contentLength != 0 {
		response += "Content-Length: " + fmt.Sprint(contentLength) + "\r\n"
	}
	response += "\r\n"
	if body != "" {
		response += body
	}
	return response
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

func readDir() (string, error) {
	directory := os.Args[2]
	if directory == "" {
		fmt.Println("Error reading file directory")
		return "", fmt.Errorf("invalid directory")
	}
	return directory, nil
}
