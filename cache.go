package materialize

import (
	"fmt"
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
	cl closer
	l  *log.Logger
}

func (w *closerWrapper) Close() {
	err := w.cl.Close()
	if err != nil {
		msg := fmt.Sprintf("failed to %T.Close: %s", w.cl, err)
		if w.l != nil {
			w.l.Print(msg)
		} else {
			log.Print(msg)
		}
	}
}

// cache caches materialized instances.
type cache struct {
	objs map[reflect.Type]reflect.Value
	c0s  []closer0
	log  *log.Logger
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
	if c0 := c.toC0(v); c0 != nil {
		c.c0s = append(c.c0s, c0)
	}
}

// closeAll closes all values which implements Close() method.
func (c *cache) closeAll() {
	for i := len(c.c0s) - 1; i >= 0; i-- {
		c.c0s[i].Close()
	}
	c.objs = map[reflect.Type]reflect.Value{}
	c.c0s = nil
}

func (c *cache) toC0(v reflect.Value) closer0 {
	typ := v.Type()
	if typ.AssignableTo(closer0Type) {
		return v.Interface().(closer0)
	}
	if typ.AssignableTo(closerType) {
		return &closerWrapper{
			cl: v.Interface().(closer),
			l:  c.log,
		}
	}
	return nil
}
