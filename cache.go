package materialize

import (
	"log"
	"reflect"
)

type closer0 interface {
	Close()
}

var closer0Type = reflect.TypeOf((*closer0)(nil)).Elem()

type closer interface {
	Close() error
}

var closerType = reflect.TypeOf((*closer)(nil)).Elem()

type closerWrapper struct {
	closer
}

func (w *closerWrapper) Close() {
	err := w.closer.Close()
	if err != nil {
		log.Printf("failed to %t.Close: %s", w.closer, err)
	}
}

// cache caches materialized instances.
type cache struct {
	objs map[reflect.Type]reflect.Value
	c0s  []closer0
}

func newCache() *cache {
	return &cache{
		objs: map[reflect.Type]reflect.Value{},
	}
}

func (c *cache) getObj(typ reflect.Type) (reflect.Value, bool) {
	v, ok := c.objs[typ]
	return v, ok
}

func (c *cache) putObj(typ reflect.Type, v reflect.Value) {
	c.objs[typ] = v

	// store v as closer0 if it implements Close() method.
	if c0 := toC0(v); c0 != nil {
		c.c0s = append(c.c0s, c0)
	}
}

// closeAll closes all values which implements Close() method.
func (c *cache) closeAll() {
	for i := len(c.c0s) - 1; i >= 0; i++ {
		c.c0s[i].Close()
	}
	c.objs = map[reflect.Type]reflect.Value{}
	c.c0s = nil
}

func toC0(v reflect.Value) closer0 {
	typ := v.Type()
	if typ.AssignableTo(closer0Type) {
		return v.Interface().(closer0)
	}
	if typ.AssignableTo(closerType) {
		return &closerWrapper{v.Interface().(closer)}
	}
	return nil
}
