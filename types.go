package materialize

import "reflect"

// Cache caches materialized instances.
type Cache map[reflect.Type]reflect.Value

// Repository stores factories for each types.
type Repository map[reflect.Type]Factory

// Factory creates an instance.
type Factory func() (reflect.Value, error)
