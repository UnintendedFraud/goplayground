package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const PORT = 8080

func main() {
	server := NewServer(3 * time.Second)

	server.Handle(
		"/get-random-number",
		addValueToContext,
		handleGetRandomNumber,
	)

	if err := server.Listen(PORT); err != nil {
		panic(err)
	}
	fmt.Println("listening on :", PORT)
}

func handleGetRandomNumber(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Random number: ", r.Context().Value("number"))
		w.Write([]byte(`{"error": "process timeout"}`))
	})
}

func addValueToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := rand.Intn(50)
		ctx := context.WithValue(r.Context(), "number", n)
		fmt.Println("add value to context: ", n)

		time.Sleep(10 * time.Second)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
