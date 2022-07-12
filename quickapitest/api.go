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
	res *http.Response,
	code int,
) T {
	t.Helper()

	path := res.Request.URL.Path
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

	got := DecodeResponse[T](t, res, http.StatusOK)
	for _, opt := range options {
		opt(t, res)
	}
	return got
}
