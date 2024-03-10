package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ServerHandler func(http.Handler) http.Handler

type Server struct {
	timeout time.Duration
}

func NewServer(timeout time.Duration) *Server {
	return &Server{
		timeout,
	}
}

func (s *Server) Handle(addr string, handlers ...ServerHandler) {
	var middlewares []ServerHandler

	if s.timeout == 0 {
		middlewares = handlers
	} else {
		middlewares = append([]ServerHandler{s.timeoutMiddleware}, handlers...)
	}

	http.Handle(addr, handleMiddlewares(middlewares))
}

func (s *Server) Listen(port int) error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		return fmt.Errorf("server broke: %s", err.Error())
	}

	return nil
}

func (s *Server) timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), s.timeout)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleMiddlewares(handlers []ServerHandler) http.Handler {
	var handler http.Handler

	for i := range handlers {
		handler = handlers[len(handlers)-1-i](handler)
	}

	return handler
}
