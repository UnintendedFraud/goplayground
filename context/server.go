package main

import (
	"fmt"
	"net/http"
)

type ServerHandler func(http.Handler) http.Handler

type Server struct {
}

func (s Server) Handle(addr string, handlers ...ServerHandler) {
	http.Handle(addr, handleMiddlewares(handlers))
}

func (s Server) Listen(port int) error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		return fmt.Errorf("server broke: %s", err.Error())
	}

	return nil
}

func handleMiddlewares(handlers []ServerHandler) http.Handler {
	var handler http.Handler

	for i := range handlers {
		handler = handlers[len(handlers)-1-i](handler)
	}

	return handler
}
