package qdump

import (
	"reflect"
	"testing"
)

func TestFill_Nil(t *testing.T) {
	want := []int{}
	got := Fill[[]int](nil)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Fill(), want=%#+v != got=%#+v", want, got)
	}
}
