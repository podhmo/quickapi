package shared

import (
	"os"
	"strconv"
)

var (
	DEBUG = false
)

func init() {
	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG")); ok {
		DEBUG = ok
	}
}
