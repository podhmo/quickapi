package qdump

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi/quickapitest"
)

func TestFillNil_Slice(t *testing.T) {
	type T = []int

	cases := []struct {
		msg   string
		want  T
		input T
	}{
		{msg: "nil", want: []int{}, input: nil},
		{msg: "empty", want: []int{}, input: []int{}},
		{msg: "values", want: []int{1, 2, 3}, input: []int{1, 2, 3}},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			want := c.want
			got := FillNil(quickapitest.NewContext(t), c.input)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFillNil_Slice2(t *testing.T) {
	type T = [][]int
	cases := []struct {
		msg   string
		want  T
		input T
	}{
		{msg: "nil", want: T{}, input: nil},
		{msg: "empty", want: T{}, input: T{}},
		{msg: "values", want: T{{1}, {}, {3}}, input: T{{1}, nil, {3}}},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			want := c.want
			got := FillNil(quickapitest.NewContext(t), c.input)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
func TestFillNil_Slice3(t *testing.T) {
	type T = [][][]int

	t.Run("nil", func(t *testing.T) {
		want := T{}
		got := FillNil[T](quickapitest.NewContext(t), nil)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("values", func(t *testing.T) {
		want := T{{{1}, {}, {2}}, {}, {{}}}
		got := FillNil(quickapitest.NewContext(t), T{{{1}, nil, {2}}, nil, {nil}})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFillNil_Map(t *testing.T) {
	type T map[string]int

	cases := []struct {
		msg   string
		want  T
		input T
	}{
		{msg: "nil", want: T{}, input: nil},
		{msg: "empty", want: T{}, input: T{}},
		{msg: "values", want: T{"foo": 0, "bar": 1}, input: T{"foo": 0, "bar": 1}},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			want := c.want
			got := FillNil(quickapitest.NewContext(t), c.input)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFillNil_Map2(t *testing.T) {
	type T map[string]map[string]int

	cases := []struct {
		msg   string
		want  T
		input T
	}{
		{msg: "nil", want: T{}, input: nil},
		{msg: "empty", want: T{}, input: T{}},
		{msg: "values",
			want:  T{"X": {"foo": 0}, "Y": {}, "Z": {"foo": 0}},
			input: T{"X": {"foo": 0}, "Y": nil, "Z": {"foo": 0}}},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			want := c.want
			got := FillNil(quickapitest.NewContext(t), c.input)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFillNil_Map3(t *testing.T) {
	type T map[string]map[string]map[string]int

	cases := []struct {
		msg   string
		want  T
		input T
	}{
		{msg: "nil", want: T{}, input: nil},
		{msg: "values",
			want: T{
				"A": {"X": {"foo": 0}, "Y": {}},
				"B": {},
			},
			input: T{
				"A": {"X": {"foo": 0}, "Y": nil},
				"B": nil,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			want := c.want
			got := FillNil(quickapitest.NewContext(t), c.input)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFillNil_Struct(t *testing.T) {
	type S struct {
		Name    string
		Friends []string
	}

	t.Run("nil", func(t *testing.T) {
		var want *S
		got := FillNil[*S](quickapitest.NewContext(t), nil)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("self-member-nil", func(t *testing.T) {
		want := S{Name: "Foo", Friends: []string{}}
		got := FillNil(quickapitest.NewContext(t), S{Name: "Foo"})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("p-self-member-nil", func(t *testing.T) {
		want := &S{Name: "Foo", Friends: []string{}}
		got := FillNil(quickapitest.NewContext(t), &S{Name: "Foo"})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFillNil_Struct_Recursive(t *testing.T) {
	type S struct {
		Name   string
		Father *S // nil -> nil
		Mother *S // nil -> nil

		Friends  []S  // nil -> []S{}
		Anothers *[]S //  nil -> nil
	}

	t.Run("nil", func(t *testing.T) {
		var want *S
		got := FillNil[*S](quickapitest.NewContext(t), nil)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("self", func(t *testing.T) {
		want := S{Name: "foo", Friends: []S{}}
		got := FillNil(quickapitest.NewContext(t), S{Name: "foo"})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("pointer-slice", func(t *testing.T) {
		want := S{Name: "foo", Friends: []S{}, Father: &S{Friends: []S{}, Name: "father"}}
		got := FillNil(quickapitest.NewContext(t), S{Name: "foo", Father: &S{Friends: nil, Name: "father"}})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("pointer-slice-slice", func(t *testing.T) {
		want := S{Name: "foo", Friends: []S{}, Father: &S{Friends: []S{{Name: "moo", Friends: []S{}}}, Name: "father"}}
		got := FillNil(quickapitest.NewContext(t), S{Name: "foo", Father: &S{Friends: []S{{Name: "moo"}}, Name: "father"}})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
		}
	})
}
