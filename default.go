package materialize

// defaultRepository is default Repository for Materializer.
var defaultRepository = Repository{}

// DefaultMaterializer provides default Materializer.
var DefaultMaterializer = New()

// Materialize gets or creates an instance of receiver's type
// with DefaultMaterializer.
func Materialize(receiver interface{}) error {
	return DefaultMaterializer.Materialize(receiver)
}

// Add adds a function as Factory with DefaultMaterializer.
func Add(fn interface{}) error {
	return DefaultMaterializer.Add(fn)
}

// MustAdd adds a function as Factory.
func MustAdd(fn interface{}) {
	DefaultMaterializer.MustAdd(fn)
}

// Close closes all cached values.
func Close() {
	DefaultMaterializer.Close()
}
