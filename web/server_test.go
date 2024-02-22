package web

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var server Server = &HTTPServer{}
	http.ListenAndServe(":8081", server)

	server.Start(":8081")
}
