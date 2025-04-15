package handler

import "net/http"

type HelloWorldHandler struct{}

func NewHelloWorldHandler() *HelloWorldHandler {
	return &HelloWorldHandler{}
}

func (h *HelloWorldHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello, World!"}`))
}
