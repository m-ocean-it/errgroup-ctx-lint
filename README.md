# errgroup-ctx-lint

## About

This linter catches cases when, within an error-group goroutine, a non-errgroup context is passed to a function or method, while there is a context, specifically attached to the said errgroup. In most cases, you want that specific context to be passed to functions invoked within the errgroup's `Go` methods.

```go
eg, egCtx := errgroup.WithContext(ctx)

eg.Go(func() error {
	return doSmth(ctx) // want "passing non-errgroup context to function withing errgroup-goroutine while there is an errgroup-context defined"

})

eg.Go(func() error {
	return doSmth(egCtx)
})

```

## TODO

- [ ] Describe the project
- [x] Respect nolint
- [ ] Add more complex test cases
- [ ] Add configuration
    - [x] custom errgroup package name
    - [ ] custom errgroup goroutine method name
- [ ] Account for wrapping of errgroup context
- [ ] suggest (and apply) fixes
- [ ] PR to golangci-lint