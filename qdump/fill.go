package qdump

import (
	"log"
	"reflect"
)

// TODO: performance up by qbind.Metadata

// Fill modifies the nil slice and maps it to an empty one, but this has side effects.
func Fill[O any](ob O) O {
	rv := reflect.ValueOf(ob)
	if rv.Kind() == reflect.Struct {
		rv = reflect.ValueOf(&ob).Elem() // for CanSet()
	}
	rv, changed := fillToplevel(rv)
	if !changed {
		return ob
	}
	return rv.Interface().(O)
}

var (
	MAX_RECURSION int = 100
)

func fillToplevel(rv reflect.Value) (ret reflect.Value, changed bool) {
	switch rv.Kind() {
	case reflect.Slice, reflect.Map:
		if sv, changed := fill(rv, 1); changed {
			return sv, true
		}
		return rv, false
	case reflect.Struct:
		return fill(rv, 1)
	case reflect.Pointer:
		_, changed := fillToplevel(rv.Elem())
		return rv, changed
	case reflect.Invalid,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Array, reflect.Chan,
		reflect.Func, reflect.Interface,
		reflect.String,
		reflect.UnsafePointer:
		return rv, false
	default:
		log.Printf("[ERROR] unsupported kind=%s, value=%v", rv.Kind(), rv)
		return rv, false
	}
}

func fill(rv reflect.Value, lv int) (ret reflect.Value, changed bool) {
	if MAX_RECURSION <= lv {
		log.Printf("[INFO] too deep lv=%d, kind=%s, value=%v", lv, rv.Kind(), rv)
		return rv, false
	}

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
				sv, subchanged := fill(rf, lv+1)
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
				sv, subchanged := fill(rf, lv+1)
				if subchanged {
					changed = true
					rv.SetMapIndex(iter.Key(), sv)
				}
			}
		}
		return rv, changed
	case reflect.Struct:
		for i, n := 0, rv.NumField(); i < n; i++ {
			rf := rv.Field(i)
			switch rf.Type().Kind() {
			case reflect.Slice, reflect.Map, reflect.Struct:
				sv, subchanged := fill(rf, lv+1)
				if subchanged {
					changed = true
					rf.Set(sv)
				}
			}
		}
		return rv, changed
	default:
		log.Printf("[ERROR] unsupported kind=%s, value=%v ...", rv.Kind(), rv)
		return rv, false
	}
}
