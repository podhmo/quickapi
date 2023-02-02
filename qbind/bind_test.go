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

type Input struct {
	Sort   string `query:"sort"`
	Pretty bool   `query:"pretty"`

	Ignored string
}

func TestBind(t *testing.T) {
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

type Input2 struct {
	*Input
	Authorization string `header:"Authorization"`
}

func TestBind2(t *testing.T) {
	metadata := qbind.Scan(func(context.Context, Input2) (interface{}, error) { return nil, nil })
	ctx := quickapitest.NewContext(t)
	req := or.Fatal(http.NewRequest("POST", "/?pretty=true&x=y&sort=-id", strings.NewReader(`{}`)))(t)
	req.Header.Set("Authorization", "bearer xxx")

	var input Input2
	if err := qbind.Bind(ctx, req, metadata, &input); err != nil {
		t.Fatalf("Bind(): unexpected error %+v", err)
	}

	want := Input2{
		Input:         &Input{Sort: "-id", Pretty: true},
		Authorization: "bearer xxx",
	}
	got := input
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Bind() mismatch (-want +got):\n%s", diff)
	}
}
