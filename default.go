package materialize

// defaultRepository is default Repository for Materializer.
var defaultRepository = &Repository{}

// DefaultMaterializer provides default Materializer.
var DefaultMaterializer = New()

// Materialize gets or creates an instance of receiver's type
// with DefaultMaterializer.
func Materialize(receiver interface{}) error {
	return DefaultMaterializer.Materialize(receiver)
}

// Add adds a function as Factory with DefaultMaterializer.
func Add(fn interface{}, tags ...string) error {
	return DefaultMaterializer.Add(fn, tags...)
}

// MustAdd adds a function as Factory.
func MustAdd(fn interface{}, tags ...string) {
	DefaultMaterializer.MustAdd(fn, tags...)
}

// CloseAll closes all cached values.
func CloseAll() {
	DefaultMaterializer.CloseAll()
}
