package quickapitest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/podhmo/quickapi/shared"
)

// DecodeResponse decodes json response
func DecodeResponse[T any](
	t *testing.T,
	method string,
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
		t.Fatalf("%s %s, status code: want=%d != got=%d", method, path, wantCode, gotCode)
	}

	var got T
	if any(got) == nil {
		t.Logf("%s %s, decode response is skipped (because nil is passed)", method, path)
		return got
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("%s %s, unexpected error (decode %T): %+v", method, path, got, err)
	}
	return got
}

// DoRequest requests and decode response
func DoRequest[T any](
	t *testing.T,
	req *http.Request,
	code int,
	handler http.Handler,
	options ...func(*testing.T, *http.Response),
) T {
	t.Helper()

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()

	got := DecodeResponse[T](t, req.Method, req.URL.Path, res, code)
	for _, opt := range options {
		opt(t, res)
	}
	return got
}

type TestLogger struct {
	T *testing.T
}

func (l *TestLogger) Printf(format string, v ...any) {
	l.T.Helper()
	l.T.Logf(format, v...)
}

func NewContext(t *testing.T) context.Context {
	return shared.SetLogger(context.Background(), &TestLogger{T: t})
}
