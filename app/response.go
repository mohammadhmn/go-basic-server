package main

import (
	"fmt"
)

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

func ResponseBuilder(statusLine string, headers map[string]string, body string) string {
	response := statusLine + "\r\n"
	for key, value := range headers {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	response += "\r\n" + body
	return response
}

func OkResponse() string {
	return OK + "\r\n\r\n"
}

func BadRequestResponse() string {
	return BAD_REQUEST + "\r\n\r\n"
}

func NotFoundResponse() string {
	return NOT_FOUND + "\r\n\r\n"
}
