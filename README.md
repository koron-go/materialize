# koron-go/materialize

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/materialize)](https://pkg.go.dev/github.com/koron-go/materialize)
[![Actions/Go](https://github.com/koron-go/materialize/workflows/Go/badge.svg)](https://github.com/koron-go/materialize/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/materialize)](https://goreportcard.com/report/github.com/koron-go/materialize)

Component's dependencies separator

## Gettings started

At first, register a factory of `*sql.DB` as component.

```go
import "github.com/koron-go/materialize"

materialize.MustAdd(func() (*sql.DB, error) {
  return sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DATASOURCE_NAME"))
})
```

Then, obtain a `*sql.DB` instance when you use.
The instance will be created automatically, and managed as singleton.

```go
var db *sql.DB
err := materialize.Materialize(&db)
if err != nil {
  return err
}
// TODO: let's work with "db".
```

All materiazlied instances which implement `Close() error` or `Close()` method,
will be closed when you call `materialize.CloseAll()`.

```go
materialize.CloseAll()
```
