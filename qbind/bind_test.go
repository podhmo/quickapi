package qbind_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/or"
	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/quickapitest"
)

func TestBind(t *testing.T) {
	type Input struct {
		Sort   string `query:"sort"`
		Pretty bool   `query:"pretty"`

		Ignored string
	}

	metadata := qbind.Scan(func(context.Context, Input) (interface{}, error) { return nil, nil })
	ctx := quickapitest.NewContext(t)
	req := or.Fatal(http.NewRequest("POST", "/?pretty=true&x=y&sort=-id", strings.NewReader(`{}`)))(t)

	var input Input
	if err := qbind.Bind(ctx, req, metadata, &input); err != nil {
		t.Fatalf("Bind(): unexpected error %+v", err)
	}

	want := Input{Sort: "-id", Pretty: true}
	got := input
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Bind() mismatch (-want +got):\n%s", diff)
	}
}
