package qdump

import (
	"log"
	"reflect"
)

// TODO: performance up by qbind.Metadata

func Fill[O any](ob O) O {
	switch rv := reflect.ValueOf(ob); rv.Kind() {
	case reflect.Slice, reflect.Map:
		if sv, changed := fill(rv); changed {
			return sv.Interface().(O)
		}
		return ob
	default:
		log.Printf("unsupported kind=%s, value=%v", rv.Kind(), rv)
		return ob
	}
}

func fill(rv reflect.Value) (ret reflect.Value, changed bool) {
	switch rv.Kind() {
	case reflect.Slice:
		if rv.IsNil() {
			return reflect.MakeSlice(rv.Type(), 0, 0), true
		}

		// slice,map,struct
		st := rv.Type().Elem()
		switch st.Kind() {
		case reflect.Slice, reflect.Map:
			for i, n := 0, rv.Len(); i < n; i++ {
				rf := rv.Index(i)
				sv, subchanged := fill(rf)
				if subchanged {
					changed = true
					rf.Set(sv)
				}
			}
		}
		return rv, changed
	case reflect.Map:
		if rv.IsNil() {
			return reflect.MakeMap(rv.Type()), true
		}

		// slice,map,struct
		st := rv.Type().Elem()
		switch st.Kind() {
		case reflect.Slice, reflect.Map:
			iter := rv.MapRange()
			for iter.Next() {
				// skip key (because: JSON's notation)
				rf := iter.Value()
				sv, subchanged := fill(rf)
				if subchanged {
					changed = true
					rv.SetMapIndex(iter.Key(), sv)
				}
			}
		}
		return rv, changed
	default:
		log.Printf("unsupported kind=%s, value=%v ...", rv.Kind(), rv)
		return rv, false
	}
}
