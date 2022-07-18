package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/podhmo/quickapi"
)

type HelloOutput struct {
	Hello string `json:"hello"`
}

func Hello(ctx context.Context, input quickapi.Empty) (output HelloOutput, err error) {
	output.Hello = "world"
	return
}

func main() {
	handler := quickapi.Lift(Hello)

	port := 8080
	log.Printf("[INFO]  listening: :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		log.Printf("[ERROR] !! %+v", err)
	}
}
