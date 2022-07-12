package qdump

import (
	"log"
	"reflect"
)

// TODO: performance up by qbind.Metadata

func Fill[O any](ob O) O {
	switch rv := reflect.ValueOf(ob); rv.Kind() {
	case reflect.Slice, reflect.Map:
		return fill(rv).Interface().(O)
	default:
		log.Printf("unsupported kind=%s, value=%v", rv.Kind(), rv)
		return ob
	}
}

func fill(rv reflect.Value) reflect.Value {
	switch rv.Kind() {
	case reflect.Slice:
		if rv.IsNil() {
			return reflect.MakeSlice(rv.Type(), 0, 0)
		}
		// slice,map,struct
		st := rv.Type().Elem()
		switch st.Kind() {
		case reflect.Slice:
			for i, n := 0, rv.Len(); i < n; i++ {
				rf := rv.Index(i)
				rf.Set(fill(rf))
			}
		}
		return rv
	case reflect.Map:
		if rv.IsNil() {
			return reflect.MakeMap(rv.Type())
		}
		return rv
	default:
		log.Printf("unsupported kind=%s, value=%v ...", rv.Kind(), rv)
		return rv
	}
}
