package materialize

import (
	"fmt"
	"reflect"
)

func (m *Materializer) materializeInterface(rv reflect.Value, typ reflect.Type, queryTags []string) error {
	realTyp, ok := m.findInterface(typ, queryTags)
	if !ok {
		return fmt.Errorf("not found assignable for: %s", typ)
	}
	return m.materializeType(rv, realTyp)
}

type matchedEntry struct {
	typ   reflect.Type
	entry *Entry
	score int
}

func (m *Materializer) findInterface(typ reflect.Type, queryTags []string) (reflect.Type, bool) {
	var me *matchedEntry
	for t, e := range m.repo {
		if !t.AssignableTo(typ) {
			continue
		}
		sc := e.Tags.score(queryTags)
		if sc >= 0 && (me == nil || me.score < sc) {
			me = &matchedEntry{
				typ:   t,
				entry: e,
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
