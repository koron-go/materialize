package materialize

import (
	"errors"
	"fmt"
	"reflect"
)

// FactoryFunc creates an instance.
type FactoryFunc func(*Context) (reflect.Value, error)

// Factory holds information for factory of a type.
type Factory struct {
	Type reflect.Type
	Func FactoryFunc
	Tags Tags
}

func (f Factory) newInstance(x *Context) (reflect.Value, error) {
	return f.Func(x.child(f.Type))
}

var (
	errType = reflect.TypeOf((*error)(nil)).Elem()
	ctxType = reflect.TypeOf((*Context)(nil))
)

func newFactory(fn interface{}) (*Factory, error) {
	rfn := reflect.ValueOf(fn)
	ft := rfn.Type()
	if ft.Kind() != reflect.Func {
		return nil, errors.New("factory should be a function")
	}

	var (
		inP  inProc
		outP outProcs
		typ  reflect.Type
	)

	switch ft.NumOut() {
	case 1:
		outP.checkLen(1)
		typ = ft.Out(0)
	case 2:
		// second of outs should be `error`.
		if !ft.Out(1).AssignableTo(errType) {
			return nil, fmt.Errorf("last of return values should be error")
		}
		outP.checkLen(2)
		outP.checkErr(1)
		typ = ft.Out(0)
	default:
		return nil, errors.New("factory should return 1 or 2 values")
	}

	switch ft.NumIn() {
	case 0:
		inP = withoutContext
	case 1:
		// type of 1st arg should be *Context
		if ft.In(0) != ctxType {
			return nil, errors.New(
				"first should be *materialize.Context if available")
		}
		inP = withContext
		outP.checkCtx()
	default:
		return nil, errors.New(
			"factory should accept no params or only *materialize.Context")
	}

	outP.checkZero()
	return &Factory{
		Type: typ,
		Func: wrapFunc(typ, rfn, inP, outP),
	}, nil
}

type inProc func(*Context) []reflect.Value

func withoutContext(*Context) []reflect.Value {
	return nil
}

func withContext(x *Context) []reflect.Value {
	return []reflect.Value{
		reflect.ValueOf(x),
	}
}

type outProc func(*Context, []reflect.Value) error

type outProcs []outProc

func (ps *outProcs) add(p ...outProc) {
	*ps = append(*ps, p...)
}

func (ps *outProcs) checkLen(expect int) {
	ps.add(func(x *Context, out []reflect.Value) error {
		n := len(out)
		if n == expect {
			return nil
		}
		panic(fmt.Sprintf("factory for %s should return %d values but %d", x.typ, expect, n))
	})
}

func (ps *outProcs) checkZero() {
	ps.add(func(x *Context, out []reflect.Value) error {
		if !out[0].IsNil() {
			return nil
		}
		return fmt.Errorf("factory for %s returned nil at 1st value", x.typ)
	})
}

func (ps *outProcs) checkErr(nerr int) {
	ps.add(func(x *Context, out []reflect.Value) error {
		rerr := out[nerr]
		if rerr.IsNil() {
			return nil
		}
		err := rerr.Interface().(error)
		return fmt.Errorf("factory for %s failed: %s", x.typ, err)
	})
}

func (ps *outProcs) checkCtx() {
	ps.add(func(x *Context, out []reflect.Value) error {
		if x.err != nil {
			return x.err
		}
		if x.val == nil {
			return nil
		}
		v := out[0].Interface()
		if v == x.val {
			return nil
		}
		panic(fmt.Sprintf("resolved value doesn't matched: resolved=%v returned=%v", x.val, v))
	})
}

// wrapFunc wraps a factory func can be used.
func wrapFunc(typ reflect.Type, fn reflect.Value, inP inProc, outP outProcs) FactoryFunc {
	return func(x *Context) (reflect.Value, error) {
		zv := reflect.Zero(typ)
		out := fn.Call(inP(x))
		for _, p := range outP {
			err := p(x, out)
			if err != nil {
				return zv, err
			}
		}
		// check context
		if x.err != nil {
			return zv, x.err
		}
		if x.val != nil {
		}
		return out[0], nil
	}
}
