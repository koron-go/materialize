package materialize

import (
	"fmt"
	"reflect"
)

// Context is a materialize context, which passed to factory as first argument.
type Context struct {
	m   *Materializer
	p   *Context
	typ reflect.Type

	val interface{}
	err error
}

func (x *Context) child(typ reflect.Type) *Context {
	return &Context{
		m:   x.m,
		p:   x,
		typ: typ,
	}
}

// Error returns last happened error if available.
func (x *Context) Error() error {
	return x.err
}

// Resolve resolves an instance temporary. This cuts circular references.
func (x *Context) Resolve(v interface{}) *Context {
	if x.val != nil {
		panic(fmt.Sprintf("have resolved already %s", x.typ))
	}
	typ := reflect.TypeOf(v)
	if typ != x.typ {
		panic(fmt.Sprintf("unmatched type, required type is %s", x.typ))
	}
	x.val = v
	return x
}

// Materialize materializes an instance with tags.
func (x *Context) Materialize(receiver interface{}, queryTags ...string) *Context {
	if x.err != nil {
		return x
	}
	x.err = x.m.materialize(x, receiver, queryTags...)
	return x
}

func (x *Context) getObj(typ reflect.Type) (reflect.Value, bool, error) {
	for x != nil {
		if x.typ == typ {
			if x.val == nil {
				return reflect.Value{}, false, fmt.Errorf("not resolved *materialize.Context for %s", x.typ)
			}
			return reflect.ValueOf(x.val), true, nil
		}
		x = x.p
	}
	return reflect.Value{}, false, nil
}
