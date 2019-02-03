package materialize

import (
	"bytes"
	"errors"
	"log"
	"testing"
)

type Res0 struct {
	sink *[]string
	id   string
}

func (r *Res0) Close() {
	*r.sink = append(*r.sink, r.id)
}

type Res1 struct {
	sink   *[]string
	id     string
	errStr string
}

func (r *Res1) Close() error {
	*r.sink = append(*r.sink, r.id)
	if r.errStr != "" {
		return errors.New(r.errStr)
	}
	return nil
}

type resFactory struct {
	r0s []*Res0
	r1s []*Res1
}

func (f *resFactory) newRes0() *Res0 {
	if len(f.r0s) == 0 {
		return nil
	}
	var v *Res0
	v, f.r0s = f.r0s[0], f.r0s[1:]
	return v
}

func (f *resFactory) newRes1() *Res1 {
	if len(f.r1s) == 0 {
		return nil
	}
	var v *Res1
	v, f.r1s = f.r1s[0], f.r1s[1:]
	return v
}

func TestCloser0(t *testing.T) {
	var sink []string
	f := &resFactory{
		r0s: []*Res0{
			{&sink, "abc"},
			{&sink, "def"},
		},
	}
	m := newTestMaterializer(t)
	m.MustAdd(f.newRes0)

	var r0 *Res0
	err := m.Materialize(&r0)
	if err != nil {
		t.Fatalf("failed to materialize Res0: %s", err)
	}
	if r0.id != "abc" {
		t.Errorf("unexpected r0.id: %s", r0.id)
	}

	if len(sink) != 0 {
		t.Fatalf("sink should be empty")
	}
	m.CloseAll()
	if len(sink) != 1 || sink[0] != "abc" {
		t.Errorf("unepected sink: %+v", sink)
	}

	var r0a *Res0
	err = m.Materialize(&r0a)
	if err != nil {
		t.Fatalf("failed to materialize 2nd Res0: %s", err)
	}
	if r0a.id != "def" {
		t.Errorf("unexpected r0a.id: %s", r0a.id)
	}
	m.CloseAll()
	if len(sink) != 2 || sink[1] != "def" {
		t.Errorf("unexpected sink after close r0a: %+v", sink)
	}

	var r0b *Res0
	err = m.Materialize(&r0b)
	if err == nil {
		t.Error("materialize r0b should be failed")
	}
	if err.Error() != "factory failed: factory for *materialize.Res0 returned nil at 1st value" {
		t.Errorf("unexpected error for r0b: %s", err)
	}
}

func TestCloser(t *testing.T) {
	var sink []string
	f := &resFactory{
		r1s: []*Res1{
			{&sink, "xxx", ""},
			{&sink, "yyy", "zzz"},
		},
	}
	bb := &bytes.Buffer{}
	m := newTestMaterializer(t).
		WithLogger(log.New(bb, "", 0))
	m.MustAdd(f.newRes1)

	var r1a *Res1
	err := m.Materialize(&r1a)
	if err != nil {
		t.Fatalf("failed to materialize Res1: %s", err)
	}
	if r1a.id != "xxx" {
		t.Errorf("unexpected r1a.id: %s", r1a.id)
	}
	m.CloseAll()
	if len(sink) != 1 || sink[0] != "xxx" {
		t.Errorf("unepected sink: %+v", sink)
	}
	if s := bb.String(); s != "" {
		t.Fatalf("unexpected logs: %q", s)
	}

	var r1b *Res1
	err = m.Materialize(&r1b)
	if err != nil {
		t.Fatalf("failed to materialize Res1: %s", err)
	}
	if r1b.id != "yyy" {
		t.Errorf("unexpected r1b.id: %q", r1b.id)
	}
	m.CloseAll()
	if len(sink) != 2 || sink[1] != "yyy" {
		t.Errorf("unepected sink: %+v", sink)
	}
	if s := bb.String(); s != "failed to *materialize.Res1.Close: zzz\n" {
		t.Fatalf("unexpected logs: %q", s)
	}
}
