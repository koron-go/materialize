package materialize

import (
	"errors"
	"testing"
)

func newTestMaterializer(t *testing.T) *Materializer {
	t.Helper()
	return New().WithRepository(Repository{})
}

type Foo struct {
	id int
}

func newFoo() *Foo {
	return &Foo{}
}

func (*Foo) Foo() {}

type Bar struct {
	id int
}

func newBar() *Bar {
	return &Bar{}
}

func (*Bar) Bar() {}

type testFactory struct {
	foo []*Foo

	bar    []*Bar
	barErr error
}

func (f *testFactory) newFoo() *Foo {
	if len(f.foo) == 0 {
		return nil
	}
	var v *Foo
	v, f.foo = f.foo[0], f.foo[1:]
	return v
}

func (f *testFactory) newBar() (*Bar, error) {
	if f.barErr != nil {
		err := f.barErr
		f.barErr = nil
		return nil, err
	}
	if len(f.bar) == 0 {
		return nil, nil
	}
	var v *Bar
	v, f.bar = f.bar[0], f.bar[1:]
	return v, nil
}

func TestMaterializeSingle(t *testing.T) {
	m := newTestMaterializer(t)
	f := &testFactory{
		foo: []*Foo{{id: 111}},
	}
	m.MustAdd(f.newFoo)

	var foo *Foo
	err := m.Materialize(&foo)
	if err != nil {
		t.Fatalf("failed to materialize *Foo: %s", err)
	}
	if foo.id != 111 {
		t.Errorf("unexpected %v: expect(id)=111", foo)
	}
}

func TestMaterializeWithError(t *testing.T) {
	m := newTestMaterializer(t)
	f := &testFactory{
		bar: []*Bar{{id: 222}},
	}
	m.MustAdd(f.newBar)

	var bar *Bar
	err := m.Materialize(&bar)
	if err != nil {
		t.Errorf("failed to materialize *Bar: %s", err)
	}
	if bar.id != 222 {
		t.Errorf("unexpected %v: expect(id)=222", bar)
	}
}

func TestMaterializeMix(t *testing.T) {
	m := newTestMaterializer(t)
	f := &testFactory{
		foo: []*Foo{{id: 333}},
		bar: []*Bar{{id: 444}},
	}
	m.MustAdd(f.newFoo)
	m.MustAdd(f.newBar)

	var err error

	var foo *Foo
	err = m.Materialize(&foo)
	if err != nil {
		t.Fatalf("failed to materialize *Foo: %s", err)
	}
	if foo.id != 333 {
		t.Errorf("unexpected %v: expect(id)=333", foo)
	}

	var bar *Bar
	err = m.Materialize(&bar)
	if err != nil {
		t.Errorf("failed to materialize *Bar: %s", err)
	}
	if bar.id != 444 {
		t.Errorf("unexpected %v: expect(id)=444", bar)
	}

	if t.Failed() {
		return
	}

	// check instance cache

	var foo2 *Foo
	err = m.Materialize(&foo2)
	if err != nil {
		t.Fatalf("failed to materialize *Foo 2nd: %s", err)
	}
	if foo2 != foo {
		t.Errorf("*Foo cache miss")
	}

	var bar2 *Bar
	err = m.Materialize(&bar2)
	if err != nil {
		t.Fatalf("failed to materialize *Bar 2nd: %s", err)
	}
	if bar2 != bar {
		t.Errorf("*Bar cache miss")
	}
}

func TestMaterializeError(t *testing.T) {
	m := newTestMaterializer(t)
	f := &testFactory{barErr: errors.New("no bars found")}
	m.MustAdd(f.newBar)
	var bar *Bar
	err := m.Materialize(&bar)
	if err == nil {
		t.Fatal("Materialize(*Bar) should faield")
	}
	if err.Error() != "factory failed: factory for *materialize.Bar failed: no bars found" {
		t.Errorf("unexpected error: %v", err)
	}
}
