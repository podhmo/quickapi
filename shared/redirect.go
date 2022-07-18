package shared

import (
	"fmt"
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

func (r redirect) Error() string {
	return fmt.Sprintf("Redirect code=%d, location=%s", r.Code, r.Location)
}

func Redirect(code int, location string) interface {
	error
	Redirector
} {
	return redirect{Code: code, Location: location}
}
