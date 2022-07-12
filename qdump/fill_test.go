package qdump

import (
	"reflect"
	"testing"
)

func TestFill_Slice(t *testing.T) {
	want := []int{}

	t.Run("nil", func(t *testing.T) {
		got := Fill[[]int](nil)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := Fill([]int{})
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
		}
	})
}
