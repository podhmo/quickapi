package qbind_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/or"
	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/quickapitest"
)

func TestBind(t *testing.T) {
	type input struct {
		Sort   string `query:"sort"`
		Pretty bool   `query:"pretty"`

		Ignored string
	}

	metadata := qbind.Scan(func(context.Context, input) (interface{}, error) { return nil, nil })
	ctx := quickapitest.NewContext(t)
	req := or.Fatal(http.NewRequest("GET", "/?pretty=true&x=y&sort=-id", nil))(t)

	got := or.Fatal(qbind.Bind[input](ctx, req, metadata))(t)
	want := input{Sort: "-id", Pretty: true}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetContext() mismatch (-want +got):\n%s", diff)
	}
}
