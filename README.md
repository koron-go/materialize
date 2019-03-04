# koron-go/materialize

[![GoDoc](https://godoc.org/github.com/koron-go/materialize?status.svg)](https://godoc.org/github.com/koron-go/materialize)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/materialize/master.svg)](https://circleci.com/gh/koron-go/materialize/tree/master)
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
