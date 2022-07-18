package shared

import (
	"fmt"
	"testing"
)

func TestStatusCodeOf(t *testing.T) {
	hmm := fmt.Errorf("hmm")

	cases := []struct {
		msg  string
		err  error
		want int
	}{
		{msg: "default", err: hmm, want: 500},
		{msg: "status-code", err: NewAPIError(hmm, 404), want: 404},
		{msg: "wrap-status-code", err: fmt.Errorf("wrap: %w", NewAPIError(hmm, 404)), want: 404},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			got := StatusCodeOf(c.err)
			if want := c.want; want != got {
				t.Errorf("StatusCodeOf(), want=%d != got=%d", want, got)
			}
		})
	}
}
