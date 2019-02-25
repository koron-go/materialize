package materialize

import "testing"

type Fooer interface {
	Foo()
}

type Barer interface {
	Bar()
}

type Bazer interface {
	Baz()
}

type FooBar struct{}

func newFooBar() *FooBar {
	return &FooBar{}
}

func (*FooBar) Foo() {}

func (*FooBar) Bar() {}

func TestMaterializeInterface(t *testing.T) {
	m := newTestMaterializer(t)
	m.MustAdd(newFoo)
	m.MustAdd(newBar)

	var f Fooer
	err := m.Materialize(&f)
	if err != nil {
		t.Fatalf("failed to materialize Fooer: %s", err)
	}
	if _, ok := f.(*Foo); !ok {
		t.Fatalf("not *Foo: %T", f)
	}

	var br Barer
	err = m.Materialize(&br)
	if err != nil {
		t.Fatalf("failed to materialize Barer: %s", err)
	}
	if _, ok := br.(*Bar); !ok {
		t.Fatalf("not *Bar: %T", br)
	}

	var bz Bazer
	err = m.Materialize(&bz)
	if err == nil {
		t.Fatalf("materialize should be failed for Bazer")
	}
	if err.Error() != "not found factory for type:materialize.Bazer tags:[]" {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestMaterializeInterfaceTags(t *testing.T) {
	m := newTestMaterializer(t)
	m.MustAdd(newFoo, "foo")
	m.MustAdd(newFooBar, "foo", "bar")

	var f1 Fooer
	err := m.Materialize(&f1, "foo")
	if err != nil {
		t.Fatalf("failed to materialize Foo(foo): %s", err)
	}
	if _, ok := f1.(*Foo); !ok {
		t.Fatalf("not *Foo: %T", f1)
	}

	var f2 Fooer
	err = m.Materialize(&f2, "foo", "bar")
	if err != nil {
		t.Fatalf("failed to materialize Foo(foo bar): %s", err)
	}
	if _, ok := f2.(*FooBar); !ok {
		t.Fatalf("not *FooBar: %T", f2)
	}
}

type strCont string

func (g strCont) Get() string {
	return string(g)
}

func newGetterFactory(s string) func() Getter {
	return func() Getter {
		return strCont(s)
	}
}

type Getter interface {
	Get() string
}

func TestAddInterface(t *testing.T) {
	m := newTestMaterializer(t)
	m.MustAdd(newGetterFactory("foo"))
	m.MustAdd(newGetterFactory("bar"), "abc")
	m.MustAdd(newGetterFactory("baz"), "xyz")

	check := func(exp string, tags ...string) {
		t.Helper()
		var g Getter
		err := m.Materialize(&g, tags...)
		if err != nil {
			t.Fatalf("failed to materialize Getter(%+v): %s", tags, err)
		}
		s := g.Get()
		if s != exp {
			t.Fatalf("unexpected Getter: %q (exp=%q)", s, exp)
		}
	}

	check("foo")
	check("bar", "abc")
	check("baz", "xyz")
}
