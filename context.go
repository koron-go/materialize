package materialize

import (
	"fmt"
	"reflect"
)

// Context is a materialize context, which passed to factory as first argument.
type Context struct {
	m *Materializer
	p *Context
	f *Factory

	val interface{}
	err error
}

func (x *Context) child(f *Factory) *Context {
	return &Context{
		m: x.m,
		p: x,
		f: f,
	}
}

// Error returns last happened error if available.
func (x *Context) Error() error {
	return x.err
}

// Resolve resolves an instance temporary. This cuts circular references.
func (x *Context) Resolve(v interface{}) *Context {
	if x.val != nil {
		panic(fmt.Sprintf("have resolved already %s", x.f.Type))
	}
	typ := reflect.TypeOf(v)
	if typ != x.f.Type {
		panic(fmt.Sprintf("unmatched type, required type is %s", x.f.Type))
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

func (x *Context) getObj(f *Factory) (reflect.Value, bool, error) {
	for x != nil {
		if x.f == f {
			if x.val == nil {
				return reflect.Value{}, false, fmt.Errorf("not resolved *materialize.Context for %s", x.f.Type)
			}
			return reflect.ValueOf(x.val), true, nil
		}
		x = x.p
	}
	return reflect.Value{}, false, nil
}

func (x *Context) typ() reflect.Type {
	return x.f.Type
}
