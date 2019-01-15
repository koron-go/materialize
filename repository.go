package materialize

import (
	"errors"
	"fmt"
	"reflect"
)

// Repository stores factories for each types.
type Repository map[reflect.Type]*Entry

var errType = reflect.TypeOf((*error)(nil)).Elem()

// Add adds a factory for a type with tags.
func (r Repository) Add(fn interface{}, tags ...string) error {
	rfn := reflect.ValueOf(fn)
	ft := rfn.Type()
	if ft.Kind() != reflect.Func {
		return errors.New("factory should be a function")
	}
	if ft.NumIn() != 0 {
		return errors.New("factory should not accept any parameters")
	}
	nout := ft.NumOut()
	if nout != 1 && nout != 2 {
		return errors.New("factory should return 1 or 2 values")
	}

	typ := ft.Out(0)
	if _, ok := r[typ]; ok {
		return fmt.Errorf("duplicated factory for %s", typ)
	}

	var fact Factory
	if nout == 2 {
		// check second outs is error.
		if !ft.Out(1).AssignableTo(errType) {
			return fmt.Errorf("last of return values should be error")
		}
		fact = newFactory2(typ, rfn)
	} else {
		fact = newFactory1(typ, rfn)
	}
	r[typ] = &Entry{
		Factory: fact,
		Tags:    newTags(tags),
	}

	return nil
}

func newFactory1(typ reflect.Type, fn reflect.Value) Factory {
	zv := reflect.Zero(typ)
	id := fmt.Sprintf("factory for %s", typ)
	return func() (reflect.Value, error) {
		out := fn.Call(nil)
		if n := len(out); n != 1 {
			return zv, fmt.Errorf("%s should retun 1 value: %d", id, n)
		}
		v := out[0]
		if v.IsNil() {
			return zv, fmt.Errorf("%s returned nil", id)
		}
		return v, nil
	}
}

func newFactory2(typ reflect.Type, fn reflect.Value) Factory {
	zv := reflect.Zero(typ)
	id := fmt.Sprintf("factory for %s", typ)
	return func() (reflect.Value, error) {
		out := fn.Call(nil)
		if n := len(out); n != 2 {
			return zv, fmt.Errorf("%s should return 2 values", id)
		}

		// check second value as error
		rerr := out[1]
		if !rerr.IsNil() {
			err := rerr.Interface().(error)
			return zv, fmt.Errorf("%s failed: %s", id, err)
		}

		v := out[0]
		if v.IsNil() {
			return zv, fmt.Errorf("%s returned nil at 1st", id)
		}
		return v, nil
	}
}
