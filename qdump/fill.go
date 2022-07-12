package qdump

import "reflect"

func Fill[O any](ob O) O {
	if val := reflect.ValueOf(ob); val.Kind() == reflect.Slice && val.IsNil() {
		ob = reflect.MakeSlice(val.Type(), 0, 0).Interface().(O)
	}
	return ob
}
