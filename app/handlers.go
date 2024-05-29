package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RequestHandler func(Request) string

func GetHandler(endpoint string) RequestHandler {
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

func HandleRoot(req Request) string {
	return OkResponse()
}

func HandleEcho(req Request) string {
	echo := strings.TrimPrefix(req.Endpoint, "/echo/")
	headers := map[string]string{
		"Content-Type":   PlainText,
		"Content-Length": fmt.Sprint(len(echo)),
	}
	encodings, ok := req.Headers["accept-encoding"]
	if ok && strings.Contains(encodings, "gzip") {
		headers["Content-Encoding"] = "gzip"
	}
	return ResponseBuilder(OK, headers, echo)
}

func HandleUserAgent(req Request) string {
	userAgent, ok := req.Headers["user-agent"]
	if !ok {
		return BadRequestResponse()
	}
	headers := map[string]string{
		"Content-Type":   PlainText,
		"Content-Length": fmt.Sprint(len(userAgent)),
	}
	return ResponseBuilder(OK, headers, userAgent)
}

func HandleFile(req Request) string {
	dir, err := readDir()
	if err != nil {
		fmt.Println(err.Error())
		return BadRequestResponse()
	}

	filename := strings.TrimPrefix(req.Endpoint, "/files/")

	if req.Method == "GET" {
		fileData, err := os.ReadFile(filepath.Join(dir, filename))
		if err != nil {
			return NotFoundResponse()
		}
		headers := map[string]string{
			"Content-Type":   OctetStream,
			"Content-Length": fmt.Sprint(len(fileData)),
		}
		return ResponseBuilder(OK, headers, string(fileData))
	} else if req.Method == "POST" {
		file, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			return NotFoundResponse()
		}
		_, err = file.WriteString(req.Body)
		if err != nil {
			return NotFoundResponse()
		}
		err = file.Close()
		if err != nil {
			return NotFoundResponse()
		}
		headers := map[string]string{
			"Content-Type": OctetStream,
		}
		return ResponseBuilder(CREATED, headers, filename)
	}

	return NotFoundResponse()
}

func readDir() (string, error) {
	directory := os.Args[2]
	if directory == "" {
		fmt.Println("Error reading file directory")
		return "", fmt.Errorf("invalid directory")
	}
	return directory, nil
}
