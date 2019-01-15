package materialize

import "reflect"

// Factory creates an instance.
type Factory func() (reflect.Value, error)

// Entry represents information of a type.
type Entry struct {
	Factory Factory
	Tags    Tags
}

func newEntry(f Factory, tags ...string) *Entry {
	return &Entry{
		Factory: f,
		Tags:    newTags(tags),
	}
}
