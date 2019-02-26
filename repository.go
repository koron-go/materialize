package materialize

import (
	"reflect"
)

// Repository stores factories for each types.
type Repository struct {
	fss map[reflect.Type]factorySet
}

// Add adds a factory for a type with tags.
func (r *Repository) Add(f *Factory) error {
	if r.fss == nil {
		r.fss = map[reflect.Type]factorySet{}
	}
	fs, ok := r.fss[f.Type]
	if !ok {
		fs = factorySet{}
		r.fss[f.Type] = fs
	}
	err := fs.add(f)
	if err != nil {
		return err
	}
	return nil
}

// Query queries a factory for type.
func (r *Repository) Query(typ reflect.Type, queryTags []string) (*Factory, bool) {
	tags := newTags(queryTags)
	mf := r.findDirect(typ, tags)
	if typ.Kind() == reflect.Interface {
		for t, fs := range r.fss {
			if t != typ && !t.AssignableTo(typ) {
				continue
			}
			mf = fs.find(mf, tags)
		}
	}
	if mf == nil {
		return nil, false
	}
	return mf.fac, true
}

// findDirect find a factory set for the type.
func (r *Repository) findDirect(typ reflect.Type, tags Tags) *matchedFactory {
	fs, ok := r.fss[typ]
	if !ok {
		return nil
	}
	return fs.find(nil, tags)
}
