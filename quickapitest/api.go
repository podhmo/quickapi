package quickapitest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// DecodeResponse decodes json response
func DecodeResponse[T any](
	t *testing.T,
	path string,
	res *http.Response,
	code int,
) T {
	t.Helper()

	defer res.Body.Close()
	if wantCode, gotCode := code, res.StatusCode; wantCode != gotCode {
		buf := new(strings.Builder)
		io.Copy(buf, res.Body)
		defer t.Logf("\tresponse: %s", buf.String())
		t.Fatalf("%q, status: want=%d != got=%d", path, wantCode, gotCode)
	}

	var got T
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("%q, unexpected error (decode %T): %+v", path, got, err)
	}
	return got
}

// DoRequest requests and decode response
func DoRequest[T any](
	t *testing.T,
	req *http.Request,
	code int,
	handler http.Handler,
	options ...func(testing.TB, *http.Response),
) T {
	t.Helper()

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()

	got := DecodeResponse[T](t, req.URL.Path, res, code)
	for _, opt := range options {
		opt(t, res)
	}
	return got
}
