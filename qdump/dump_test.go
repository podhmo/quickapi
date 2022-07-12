package qdump

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConnectionIsClosed(t *testing.T) {
	type T = interface{}

	cases := []struct {
		msg      string
		wantCode int
		err      error
		request  *http.Request
	}{
		{
			msg: "connection-is-closed-in-client", wantCode: 499, err: context.Canceled, request: func() *http.Request {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // canceled context
				return httptest.NewRequest("GET", "/", nil).WithContext(ctx)
			}(),
		},
		{
			msg: "internal-timeout", wantCode: 500, err: context.Canceled, request: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil).WithContext(context.Background())
			}(),
		},
	}

	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			rec := httptest.NewRecorder()
			Dump[T](rec, c.request, nil, context.Canceled)

			if want, got := c.wantCode, rec.Result().StatusCode; want != got {
				t.Errorf("status-code in Dump(), want=%d != got=%d", want, got)
			}
		})
	}
}
