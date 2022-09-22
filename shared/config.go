package shared

import (
	"os"
	"strconv"
)

var (
	DEBUG = false
	INFO  = false
)

func init() {
	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG")); ok {
		DEBUG = ok
		INFO = ok
	}
	if ok, _ := strconv.ParseBool(os.Getenv("INFO")); ok {
		INFO = ok
	}
}
