package pathutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNormalizeTemplatedPath(t *testing.T) {
	cases := []struct {
		path       string
		vars       map[string]string
		normalized string
	}{
		{path: "/foo", normalized: "/foo", vars: nil},
		{path: "/foo/{id}", normalized: "/foo/{id}", vars: map[string]string{"id": ""}},
		{path: "/foo/{ id}", normalized: "/foo/{id}", vars: map[string]string{"id": ""}},
		{path: "/foo/{id }", normalized: "/foo/{id}", vars: map[string]string{"id": ""}},
		{path: "/foo/{ id }", normalized: "/foo/{id}", vars: map[string]string{"id": ""}},
		{path: "/foo/{foo_id}/bar/{bar_id}", normalized: "/foo/{foo_id}/bar/{bar_id}", vars: map[string]string{"foo_id": "", "bar_id": ""}},
		{path: "/foo/{foo_id}/articles/{rid:^[0-9]{5,6}}", normalized: "/foo/{foo_id}/articles/{rid}", vars: map[string]string{"foo_id": "", "rid": "^[0-9]{5,6}"}},
		{path: "/foo/{id:^(number)?:?[0-9]+}", normalized: "/foo/{id}", vars: map[string]string{"id": "^(number)?:?[0-9]+"}},
		// RFC6570 foo*
		// - https://github.com/OAI/OpenAPI-Specification/issues/291
		// - https://github.com/OAI/OpenAPI-Specification/issues/892

		// https://www.rfc-editor.org/rfc/rfc6570#section-3.2.1 	with count := ("one", "two", "three")
		// - {count}            one,two,three
		// - {count*}           one,two,three
		// - {/count}           /one,two,three
		// - {/count*}          /one/two/three
		{path: "/api/metadata/*", normalized: "/api/metadata/{STAR*}", vars: map[string]string{"STAR*": ""}},
		{path: "/version/{id}/api/metadata/*", normalized: "/version/{id}/api/metadata/{STAR*}", vars: map[string]string{"id": "", "STAR*": ""}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.path, func(t *testing.T) {
			normalized, _, vars := NormalizeTemplatedPath(c.path)
			{
				type ref struct{ Vars map[string]string }
				want := c.vars
				got := vars
				if diff := cmp.Diff(ref{want}, ref{got}); diff != "" {
					t.Errorf("normalizeTemplatePath(), vars mismatch (-want +got):\n%s", diff)
				}
			}
			{
				want := c.normalized
				got := normalized
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("normalizeTemplatePath(), normalied path mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
