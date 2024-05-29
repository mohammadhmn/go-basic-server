package main

import (
	"bufio"
	"strings"
)

type Request struct {
	Method   string
	Endpoint string
	Headers  map[string]string
	Body     string
}

func ParseRequest(requestString string) Request {
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
