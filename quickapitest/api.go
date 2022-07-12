package quickapitest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// DecodeResponse decodes json response
func DecodeResponse[T any](t *testing.T, path string, res *http.Response, code int) T {
	t.Helper()

	defer res.Body.Close()
	var got T
	if wantCode, gotCode := code, res.StatusCode; wantCode != gotCode {
		buf := new(strings.Builder)
		io.Copy(buf, res.Body)
		defer t.Logf("\tresponse: %s", buf.String())
		t.Fatalf("%q, status: want=%d != got=%d", path, wantCode, gotCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("%q, unexpected error (decode %T): %+v", path, got, err)
	}
	return got
}
