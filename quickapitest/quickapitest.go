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
	"golang.org/x/exp/slog"
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
		defer t.Logf("\tbody: %s", buf.String())
		t.Fatalf("response: %-7s %s -- status code: want=%d != got=%d", method, path, wantCode, gotCode)
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

// DoHandler requests and decode response
func DoHandler[T any](
	t *testing.T,
	handler http.Handler,
	req *http.Request,
	code int,
	options ...func(*testing.T, *http.Response),
) T {
	t.Helper()

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if l := shared.GetLoggerOrNil(req.Context()); l == nil {
		req = req.WithContext(shared.SetLogger(req.Context(), NewTestLogger(t)))
	}
	t.Logf("request : %-7s %s -- with-body?=%4v", req.Method, req.URL, req.Body != nil && req.Body != http.NoBody)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()

	got := DecodeResponse[T](t, req.Method, req.URL.Path, res, code)
	for _, opt := range options {
		opt(t, res)
	}
	return got
}

// TODO: colorful output (pretty output)

type TestWriter struct {
	T *testing.T
}

func (l *TestWriter) Write(b []byte) (int, error) {
	l.T.Helper()
	l.T.Logf(string(b))
	return len(b), nil
}

func NewContext(t *testing.T) context.Context {
	return shared.SetLogger(context.Background(), NewTestLogger(t))
}

func NewTestLogger(t *testing.T) *slog.Logger {
	// remove fields {:time:, :source:}
	return slog.New(slog.NewJSONHandler(&TestWriter{T: t}, &slog.HandlerOptions{
		AddSource: shared.DEBUG,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
}
