package qdump

import (
	"net/http"
)

type Redirector interface {
	Redirect(http.ResponseWriter, *http.Request)
}

type redirect struct {
	Location string
	Code     int
}

func (r redirect) Redirect(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Location", r.Location)
	w.WriteHeader(r.Code)
}

func Redirect(code int, location string) Redirector {
	return redirect{Code: code, Location: location}
}
