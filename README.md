# quickapi (WIP)

not fast implementation, but quick web API development

> **Warning**
> 🚧 This package is under construction, so all examples may not work correctly.

We need the type just like this.

```go
type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)
```

and lift it.

```go
func ListTodo(context.Context, quickapi.Empty) ([]Todo, error) {
    return []Todo{{Title: "foo"}}, nil
}

r.Get("/todos", quickapi.Lift(ListTodo))
```

## experimental openapi support

:warning: this is experimental feature, using quickapi/qopenapi/define package can be able to define openapi doc with define router.


## how to use

see [examples](_examples)
or https://github.com/podhmo/gtasks-server

