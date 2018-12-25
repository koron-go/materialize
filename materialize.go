package materialize

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Materializer manages materialize instances.
type Materializer struct {
	mu    sync.Mutex
	cache Cache

	Repository Repository
}

// New creates a Materializer.
func New() *Materializer {
	return &Materializer{
		cache: Cache{},
	}
}

// Materialize gets or creates an instance of receiver's type.
func (m *Materializer) Materialize(receiver interface{}) error {
	rv := reflect.ValueOf(receiver)
	rt := rv.Type()
	if rt.Kind() != reflect.Ptr {
		return errors.New("receiver should be a pointer")
	}
	typ := rt.Elem()

	m.mu.Lock()
	defer m.mu.Unlock()

	if v, ok := m.cache[typ]; ok {
		rv.Elem().Set(v)
		return nil
	}

	f, ok := m.repos()[typ]
	if !ok {
		return fmt.Errorf("not found factories for: %s", typ)
	}
	v, err := f()
	if err != nil {
		return fmt.Errorf("factory failed: %v", err)
	}
	m.cache[typ] = v
	rv.Elem().Set(v)

	return nil
}

func (m *Materializer) repos() Repository {
	if m.Repository != nil {
		return m.Repository
	}
	return DefaultRepository
}

func (m *Materializer) addFactory(typ reflect.Type, f Factory) {
	m.repos()[typ] = f
}

// MustAdd adds a function as Factory.
func (m *Materializer) MustAdd(fn interface{}) {
	err := m.Add(fn)
	if err != nil {
		panic(err)
	}
}

var errType = reflect.TypeOf((*error)(nil)).Elem()

// Add adds a function as Factory.
func (m *Materializer) Add(fn interface{}) error {
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

	m.mu.Lock()
	defer m.mu.Unlock()

	typ := ft.Out(0)
	if _, ok := m.repos()[typ]; ok {
		return fmt.Errorf("duplicated factory for %s", typ)
	}

	if nout == 2 {
		// check second outs is error.
		if !ft.Out(1).AssignableTo(errType) {
			return fmt.Errorf("last of return values should be error")
		}
		m.addFactory(typ, newFactory2(typ, rfn))
	} else {
		m.addFactory(typ, newFactory1(typ, rfn))
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
