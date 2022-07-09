package qbind

import (
	"encoding/json"
	"log"
	"net/http"
)

func Bind[I any](req *http.Request, isEmpty bool) (I, error) {
	var input I
	if !isEmpty {
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			log.Printf("[ERROR] unexpected error (json.Decode): %+v", err) // TODO: structured logging
			return input, err
		}
	}
	return input, nil
}
