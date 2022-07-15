package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/podhmo/quickapi"
)

func main() {
	r := quickapi.DefaultRouter()
	r.Get("/", quickapi.Lift(Hello))

	port := 8080
	log.Printf("[Info]  listening: :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
}

type HelloOutput struct {
	Hello string `json:"hello"`
}

func Hello(ctx context.Context, input quickapi.Empty) (output HelloOutput, err error) {
	output.Hello = "world"
	return
}
