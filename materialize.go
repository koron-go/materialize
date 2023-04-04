package materialize

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

// Materializer manages materialize instances.
type Materializer struct {
	mu    sync.Mutex
	cache *cache
	repo  *Repository
	log   *log.Logger

	currRootX *Context
}

// New creates a Materializer.
func New() *Materializer {
	return &Materializer{
		cache: newCache(),
	}
}

// WithRepository replaces a Repository.
func (m *Materializer) WithRepository(r *Repository) *Materializer {
	m.repo = r
	return m
}

// WithLogger replaces a *log.Logger.
func (m *Materializer) WithLogger(l *log.Logger) *Materializer {
	m.log = l
	m.cache.log = l
	return m
}

// unused, disabled
//func (m *Materializer) logf(format string, args ...interface{}) {
//	if m.log == nil {
//		log.Printf(format, args...)
//		return
//	}
//	m.log.Printf(format, args...)
//}

// Materialize gets or creates an instance of receiver's type.
func (m *Materializer) Materialize(receiver interface{}, queryTags ...string) error {
	if m.currRootX != nil {
		return ErrorBusy
	}
	m.mu.Lock()
	defer func() {
		m.currRootX = nil
		m.mu.Unlock()
	}()
	x := &Context{m: m}
	m.currRootX = x
	return m.materialize(x, receiver, queryTags)
}

func (m *Materializer) materialize(x *Context, receiver interface{}, queryTags []string) error {
	rv := reflect.ValueOf(receiver)
	if rv.Kind() != reflect.Ptr {
		return ErrorReceiverType
	}
	typ := rv.Type().Elem()

	switch typ.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Ptr, reflect.String:
		return m.materialize0(x, rv, typ, queryTags)
	default:
		return fmt.Errorf("unsupported type:%s kind:%s", typ, typ.Kind())
	}
}

// materialize0 materializes an object for the factory.
func (m *Materializer) materialize0(x *Context, rv reflect.Value, typ reflect.Type, queryTags []string) error {
	f, ok := m.getRepo().Query(typ, queryTags)
	if !ok {
		return fmt.Errorf("not found factory for type:%s tags:%+v", typ, queryTags)
	}

	v0, ok, err := x.getObj(f)
	if err != nil {
		return err
	} else if ok {
		rv.Elem().Set(v0)
		return nil
	}

	if v, ok := m.cache.getObj(f); ok {
		rv.Elem().Set(v)
		return nil
	}

	v, err := f.newInstance(x)
	if err != nil {
		return fmt.Errorf("factory failed: %w", err)
	}
	m.cache.putObj(f, v)
	rv.Elem().Set(v)

	return nil
}

func (m *Materializer) getRepo() *Repository {
	if m.repo != nil {
		return m.repo
	}
	return defaultRepository
}

// MustAdd adds a function as Factory.
func (m *Materializer) MustAdd(fn interface{}, tags ...string) *Materializer {
	err := m.Add(fn, tags...)
	if err != nil {
		panic(err)
	}
	return m
}

// Add adds a function as Factory.
func (m *Materializer) Add(fn interface{}, tags ...string) error {
	f, err := newFactory(fn, tags)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	err = m.getRepo().Add(f)
	if err != nil {
		return err
	}
	return nil
}

// CloseAll closes all values which implements Close() method, and clear value
// cache.
func (m *Materializer) CloseAll() {
	m.mu.Lock()
	m.cache.closeAll()
	m.mu.Unlock()
}
