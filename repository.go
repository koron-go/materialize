package materialize

import (
	"fmt"
	"reflect"
)

// Repository stores factories for each types.
type Repository map[reflect.Type]*Factory

// Add adds a factory for a type with tags.
func (r Repository) Add(fn interface{}, tags ...string) error {
	f, err := newFactory(fn)
	if err != nil {
		return err
	}
	if _, ok := r[f.Type]; ok {
		return fmt.Errorf("duplicated factory for %s", f.Type)
	}
	f.Tags = newTags(tags)
	r[f.Type] = f
	return nil
}
