# quickapi (WIP)

not fast, but quick

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

## how to use

see [examples](_examples)