package materialize

import (
	"fmt"
	"reflect"
)

// Repository stores factories for each types.
type Repository struct {
	facs map[reflect.Type]*Factory
}

// Add adds a factory for a type with tags.
func (r *Repository) Add(f *Factory, tags ...string) error {
	if _, ok := r.facs[f.Type]; ok {
		return fmt.Errorf("duplicated factory for %s", f.Type)
	}
	f.Tags = newTags(tags)
	if r.facs == nil {
		r.facs = map[reflect.Type]*Factory{}
	}
	r.facs[f.Type] = f
	return nil
}

// Get gets a factory for type.
func (r *Repository) Get(typ reflect.Type) (*Factory, bool) {
	f, ok := r.facs[typ]
	return f, ok
}

type matchedFactory struct {
	fac *Factory
	sc  int
}

func (r *Repository) findInterface(typ reflect.Type, tags []string) (*Factory, bool) {
	var mf *matchedFactory
	for t, f := range r.facs {
		if !t.AssignableTo(typ) {
			continue
		}
		sc := f.Tags.score(tags)
		if sc >= 0 && (mf == nil || mf.sc < sc) {
			mf = &matchedFactory{
				fac: f,
				sc:  sc,
			}
			// FIXME: log matchedFactory
		}
	}
	if mf == nil {
		return nil, false
	}
	return mf.fac, true
}
