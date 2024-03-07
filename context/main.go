package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

const PORT = 8080

func main() {
	server := Server{}

	server.Handle(
		"/get-random-number",
		printHello,
		handleGetRandomNumber,
	)

	if err := server.Listen(PORT); err != nil {
		panic(err)
	}
	fmt.Println("listening on :", PORT)
}

func handleGetRandomNumber(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Random number: ", rand.Int())
	})
}

func printHello(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello")

		next.ServeHTTP(w, r)
	})
}
