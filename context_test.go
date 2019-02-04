package materialize

import "testing"

func TestContext_WithError(t *testing.T) {
	m := newTestMaterializer(t)
	var exp *Foo
	m.MustAdd(func(x *Context) (*Foo, error) {
		v := &Foo{id: 1234}
		x.Resolve(v)
		exp = v
		return v, nil
	})
	var foo *Foo
	err := m.Materialize(&foo)
	if err != nil {
		t.Fatalf("failed to Materialize: %s", err)
	}
	if exp == nil {
		t.Fatal("factory didn't be called")
	}
	if foo != exp {
		t.Fatal("materialized Foo is not matched")
	}
}

func TestContext_WithoutError(t *testing.T) {
	m := newTestMaterializer(t)
	var exp *Foo
	m.MustAdd(func(x *Context) *Foo {
		v := &Foo{id: 1234}
		x.Resolve(v)
		exp = v
		return v
	})
	var foo *Foo
	err := m.Materialize(&foo)
	if err != nil {
		t.Fatalf("failed to Materialize: %s", err)
	}
	if foo != exp {
		t.Fatal("materialized Foo is not matched")
	}
}

func TestContext_WithoutResolve(t *testing.T) {
	m := newTestMaterializer(t)
	var exp *Foo
	m.MustAdd(func(x *Context) *Foo {
		v := &Foo{id: 1234}
		exp = v
		return v
	})
	var foo *Foo
	err := m.Materialize(&foo)
	if err != nil {
		t.Fatalf("failed to Materialize: %s", err)
	}
	if foo != exp {
		t.Fatal("materialized Foo is not matched")
	}
}

type FooX struct {
	foo *Foo
	bar *Bar
}

func TestContext_Materialize(t *testing.T) {
	m := newTestMaterializer(t)
	var (
		z0 *FooX
		z1 *Foo
		z2 *Bar
	)
	m.MustAdd(func(x *Context) (*FooX, error) {
		z0 = &FooX{}
		x.Resolve(z0).Materialize(&z0.foo).Materialize(&z0.bar)
		return z0, nil
	}).MustAdd(func(x *Context) (*Foo, error) {
		z1 = &Foo{id: 1234}
		return z1, nil
	}).MustAdd(func(x *Context) (*Bar, error) {
		z2 = &Bar{id: 9876}
		return z2, nil
	})

	var fooX *FooX
	err := m.Materialize(&fooX)
	if err != nil {
		t.Fatalf("failed to Materialize: %s", err)
	}

	if fooX != z0 {
		t.Fatalf("materialized FooX unmatched: %+v %+v", fooX, z0)
	}
	if fooX.foo != z1 {
		t.Fatalf("materialized FooX.foo unmatched: %+v %+v", fooX.foo, z1)
	}
	if fooX.bar != z2 {
		t.Fatalf("materialized FooX.bar unmatched: %+v %+v", fooX.bar, z1)
	}
	if fooX.foo.id != 1234 {
		t.Fatalf("unexpected fooX foo: %+v", fooX.foo)
	}
	if fooX.bar.id != 9876 {
		t.Fatalf("unexpected fooX bar: %+v", fooX.bar)
	}
}

type Circular1A struct {
	b *Circular1B
}

type Circular1B struct {
	a *Circular1A
}

func TestContext_Circular(t *testing.T) {
	m := newTestMaterializer(t)
	var (
		zA *Circular1A
		zB *Circular1B
	)
	m.MustAdd(func (x *Context) *Circular1A {
		zA = &Circular1A{}
		x.Resolve(zA).Materialize(&zA.b)
		return zA
	}).MustAdd(func(x *Context) *Circular1B {
		zB = &Circular1B{}
		x.Resolve(zB).Materialize(&zB.a)
		return zB
	})

	var rA *Circular1A
	err := m.Materialize(&rA)
	if err != nil {
		t.Fatalf("failed to Materialize %T: %s", rA, err)
	}
	if rA != zA {
		t.Fatalf("materialized rA unmatched: %+v %+v", rA, zA)
	}
	if rA.b != zB {
		t.Fatal("materialized rA.b unmatched")
	}
	if rA.b.a != zA {
		t.Fatal("materialized rA.b.a unmatched")
	}
	// completed circular

	var rB *Circular1B
	err = m.Materialize(&rB)
	if err != nil {
		t.Fatalf("failed to Materialize %T: %s", rB, err)
	}
	if rB != zB {
		t.Fatalf("materialized rB unmatched: %+v %+v", rB, zB)
	}
	if rB.a != zA {
		t.Fatal("materialized rB.b unmatched")
	}
	if rB.a.b != zB {
		t.Fatal("materialized rB.b.a unmatched")
	}
}
