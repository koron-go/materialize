package materialize

import "reflect"

// Repository stores factories for each types.
type Repository map[reflect.Type]Factory

// Factory creates an instance.
type Factory func() (reflect.Value, error)
