package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const PORT = 8080

func main() {
	server := NewServer(1 * time.Second)

	server.Handle(
		"/get-value",
		addValueToContext,
		handleGetValue,
	)

	if err := server.Listen(PORT); err != nil {
		panic(err)
	}
	fmt.Println("listening on :", PORT)
}

func handleGetValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("starting handleGetValue")
		defer fmt.Println("finished executing handleGetNumber")

		message, err := executeHandleGetValue(r.Context(), r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("\nError happened: %s", err.Error())))
			fmt.Println("Error happened: ", err)
			return
		}

		w.Write([]byte(message))

	})
}

func executeHandleGetValue(ctx context.Context, r *http.Request) (string, error) {
	chanErr := make(chan error, 1)

	go func() {
		fmt.Println("executeHandleGetValue started")

		time.Sleep(2 * time.Second)

		// simulate errors
		shouldError := true

		if shouldError {
			chanErr <- fmt.Errorf("an error happened during goroutine")
		} else {
			chanErr <- nil
		}

		fmt.Println("executeHandleGetValue ended")
	}()

	select {
	case <-ctx.Done():
		<-chanErr
		return "", ctx.Err()
	case err := <-chanErr:
		return "finished successfully", err
	}
}

func addValueToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add simple values
		ctx1 := context.WithValue(r.Context(), "number", 420)
		ctx2 := context.WithValue(ctx1, "sad_message", "RIP Toriyama :(")

		// Add a more complex one
		cs := ComplexStruct{
			question: "What is your favourite Dragon Ball character?",
			possibleAnswers: []string{
				"Goku",
				"Gohan",
				"Vegeta",
				"You get the idea zz",
			},
		}

		ctx3 := context.WithValue(ctx2, "complex_struct", cs)

		next.ServeHTTP(w, r.WithContext(ctx3))
	})
}

type ComplexStruct struct {
	question        string
	possibleAnswers []string
}
