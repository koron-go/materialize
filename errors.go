package materialize

import "errors"

var (
	// ErrorBusy is an error when called Materializer.Materialize in materialization context.
	ErrorBusy = errors.New("busy or recursive materialization, try materialize.Context#Materialize() instead")

	// ErrorReceiverType shows receiver is not a pointer type.
	ErrorReceiverType = errors.New("receiver should be a pointer")

	// ErrorFactoryType shows a factory is not expected type (function).
	ErrorFactoryType = errors.New("factory should be a function")

	// ErrorFactoryRetun shows a factory has unexpected number of return values.
	ErrorFactoryRetun = errors.New("factory should return 1 or 2 values")

	// ErrorFactoryFirstArg shows a factory should have *materialize.Context as 1st argument.
	ErrorFactoryFirstArg = errors.New("first should be *materialize.Context if available")

	// ErrorFactoryArgsRule shows a factory should have no arguments or just one argument (*materialize.Context).
	ErrorFactoryArgsRule = errors.New("factory should accept no params or only *materialize.Context")
)
