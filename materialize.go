package materialize

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
)

// Materializer manages materialize instances.
type Materializer struct {
	mu    sync.Mutex
	cache *cache
	repo  Repository
	log   *log.Logger
}

// New creates a Materializer.
func New() *Materializer {
	return &Materializer{
		cache: newCache(),
	}
}

// WithRepository replaces a Repository.
func (m *Materializer) WithRepository(r Repository) *Materializer {
	m.repo = r
	return m
}

// WithLogger replaces a *log.Logger.
func (m *Materializer) WithLogger(l *log.Logger) *Materializer {
	m.log = l
	m.cache.log = l
	return m
}

func (m *Materializer) logf(format string, args ...interface{}) {
	if m.log == nil {
		log.Printf(format, args...)
		return
	}
	m.log.Printf(format, args...)
}

// Materialize gets or creates an instance of receiver's type.
func (m *Materializer) Materialize(receiver interface{}, queryTags ...string) error {
	rv := reflect.ValueOf(receiver)
	if rv.Kind() != reflect.Ptr {
		return errors.New("receiver should be a pointer")
	}
	typ := rv.Type().Elem()

	m.mu.Lock()
	defer m.mu.Unlock()

	switch typ.Kind() {
	case reflect.Ptr:
		return m.materializeType(rv, typ)
	case reflect.Interface:
		return m.materializeInterface(rv, typ, queryTags)
	default:
		return fmt.Errorf("unsupported type: %s (%s)", typ, typ.Kind())
	}
}

func (m *Materializer) materializeType(rv reflect.Value, typ reflect.Type) error {
	if v, ok := m.cache.getObj(typ); ok {
		rv.Elem().Set(v)
		return nil
	}

	e, ok := m.getRepo()[typ]
	if !ok {
		return fmt.Errorf("not found factories for: %s", typ)
	}
	v, err := e.Factory()
	if err != nil {
		return fmt.Errorf("factory failed: %v", err)
	}
	m.cache.putObj(typ, v)
	rv.Elem().Set(v)

	return nil
}

func (m *Materializer) getRepo() Repository {
	if m.repo != nil {
		return m.repo
	}
	return defaultRepository
}

// MustAdd adds a function as Factory.
func (m *Materializer) MustAdd(fn interface{}, tags ...string) {
	err := m.Add(fn, tags...)
	if err != nil {
		panic(err)
	}
}

// Add adds a function as Factory.
func (m *Materializer) Add(fn interface{}, tags ...string) error {
	m.mu.Lock()
	err := m.getRepo().Add(fn, tags...)
	m.mu.Unlock()
	return err
}

// CloseAll closes all values which implements Close() method, and clear value
// cache.
func (m *Materializer) CloseAll() {
	m.mu.Lock()
	m.cache.closeAll()
	m.mu.Unlock()
}
