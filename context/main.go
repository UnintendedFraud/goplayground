package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const PORT = 8080

func main() {
	server := &Server{}

	server.Handle("/get-value", addValueToContext, handleGetValue)

	if err := server.Listen(PORT); err != nil {
		panic(err)
	}
	fmt.Println("listening on :", PORT)
}

func handleGetValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("handleGetValue started")
		defer log.Println("handleGetNumber ended")

		ctx, cancelCtx := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancelCtx()

		err := someLongAction(ctx)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("\n%s - Error happened: %s", t(), err.Error())))
			log.Println("Error happened: ", err)
			return
		}

		w.Write([]byte(fmt.Sprintf("%s -- operation finished successfully", t())))
	})
}

func someLongAction(ctx context.Context) error {
	log.Println("someLongAction started")
	defer log.Println("someLongAction ended")

	select {
	case err := <-simulatingOperation():
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func simulatingOperation() chan error {
	log.Println("simlatingOperation started")
	defer log.Println("simlatingOperation ended")

	chanErr := make(chan error, 1)

	go func() {
		log.Println("goroutine in simulatingOperation started")
		defer log.Println("goroutine in simulatingOperation ended")

		time.Sleep(5 * time.Second)
		chanErr <- fmt.Errorf("something terrible happened, PLEASE HELP!")
	}()

	return chanErr
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

func t() string {
	return time.Now().Format("15:06:05")
}

type ComplexStruct struct {
	question        string
	possibleAnswers []string
}
