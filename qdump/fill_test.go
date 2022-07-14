package qdump

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFill_Slice(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		want := []int{}
		got := Fill[[]int](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		want := []int{}
		got := Fill(want)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("values", func(t *testing.T) {
		want := []int{1, 2, 3}
		got := Fill(want)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})
}

func TestFill_Slice2(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		want := [][]int{}
		got := Fill[[][]int](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("nil-values2", func(t *testing.T) {
		want := [][]int{{1}, {}, {3}}
		got := Fill([][]int{{1}, nil, {3}})
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("nil-values3", func(t *testing.T) {
		want := [][][]int{{{1}, {}, {2}}, {}, {{}}}
		got := Fill([][][]int{{{1}, nil, {2}}, nil, {nil}})
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("nil-values3-nil", func(t *testing.T) {
		want := [][][]int{}
		got := Fill[[][][]int](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})
}

func TestFill_Map(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		want := map[string]int{}
		got := Fill[map[string]int](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		want := map[string]int{}
		got := Fill(want)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("values", func(t *testing.T) {
		want := map[string]int{"foo": 0, "bar": 1}
		got := Fill(want)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})
}

func TestFill_Map2(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		want := map[string]map[string]int{}
		got := Fill[map[string]map[string]int](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("nil-values2", func(t *testing.T) {
		want := map[string]map[string]int{
			"X": {"foo": 0},
			"Y": {},
			"Z": {"foo": 0},
		}
		got := Fill(map[string]map[string]int{
			"X": {"foo": 0},
			"Y": nil,
			"Z": {"foo": 0},
		})
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("nil-values3", func(t *testing.T) {
		want := map[string]map[string]map[string]int{
			"A": {"X": {"foo": 0}, "Y": {}},
			"B": {},
		}
		got := Fill(map[string]map[string]map[string]int{
			"A": {"X": {"foo": 0}, "Y": nil},
			"B": nil,
		})
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("nil-values3-nil", func(t *testing.T) {
		want := map[string]map[string]map[string]int{}
		got := Fill[map[string]map[string]map[string]int](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})
}

func TestFill_Struct(t *testing.T) {
	type S struct {
		Name    string
		Friends []string
	}

	t.Run("nil", func(t *testing.T) {
		var want *S
		got := Fill[*S](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("member-nil", func(t *testing.T) {
		want := S{Name: "Foo", Friends: []string{}}
		got := Fill(S{Name: "Foo"})
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("p-member-nil", func(t *testing.T) {
		want := &S{Name: "Foo", Friends: []string{}}
		got := Fill(&S{Name: "Foo"})
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})
}

func TestFill_Struct_Recursive(t *testing.T) {
	type S struct {
		Name   string
		Father *S
		Mother *S

		Friends  []S
		Anothers *[]S
	}

	t.Run("nil", func(t *testing.T) {
		var want *S
		got := Fill[*S](nil)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fill() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("self", func(t *testing.T) {
		want := S{Name: "foo", Friends: []S{}}
		got := Fill(S{Name: "foo"})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fill() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("pointer-slice", func(t *testing.T) {
		want := S{Name: "foo", Friends: []S{}, Father: &S{Friends: []S{}, Name: "father"}}
		got := Fill(S{Name: "foo", Father: &S{Friends: nil, Name: "father"}})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fill() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("pointer-slice-slice", func(t *testing.T) {
		want := S{Name: "foo", Friends: []S{}, Father: &S{Friends: []S{{Name: "moo", Friends: []S{}}}, Name: "father"}}
		got := Fill(S{Name: "foo", Father: &S{Friends: []S{{Name: "moo"}}, Name: "father"}})
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fill() mismatch (-want +got):\n%s", diff)
		}
	})
}
