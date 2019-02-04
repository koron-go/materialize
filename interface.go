package materialize

import (
	"fmt"
	"reflect"
)

func (m *Materializer) materializeInterface(x *Context, rv reflect.Value, typ reflect.Type, queryTags []string) error {
	realTyp, ok := m.findInterface(typ, queryTags)
	if !ok {
		return fmt.Errorf("not found assignable for: %s", typ)
	}
	return m.materializeType(x, rv, realTyp)
}

type matchedEntry struct {
	typ   reflect.Type
	entry *Factory
	score int
}

func (m *Materializer) findInterface(typ reflect.Type, queryTags []string) (reflect.Type, bool) {
	var me *matchedEntry
	for t, f := range m.repo {
		if !t.AssignableTo(typ) {
			continue
		}
		sc := f.Tags.score(queryTags)
		if sc >= 0 && (me == nil || me.score < sc) {
			me = &matchedEntry{
				typ:   t,
				entry: f,
				score: sc,
			}
			// FIXME: log matchedEntry
		}
	}
	if me == nil {
		return nil, false
	}
	return me.typ, true
}
