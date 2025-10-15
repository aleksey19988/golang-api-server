package http

import (
	"io"
	"net/http"
)

func PingHandler(w http.ResponseWriter, req *http.Request) {
	_, err := io.WriteString(w, "pong\n")
	if err != nil {
		return
	}
}
