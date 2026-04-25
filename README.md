# errgroup-ctx-lint


## About

This linter catches cases when, within an error-group goroutine, a non-errgroup context is passed to a function or method, while there is a context, specifically attached to the said errgroup. In most cases, you want that specific context to be passed to functions invoked within the errgroup's `Go` methods.

```go
eg, egCtx := errgroup.WithContext(ctx)

eg.Go(func() error {
	return doSmth(ctx) // want `errgroup callback should probably not reference outer context "ctx", use the errgroup-derived context "egCtx"`
})

eg.Go(func() error {
	return doSmth(egCtx) // Correctly uses the context returned by "errgroup.WithContext"
})

eg.TryTo(func() error {
	return doSmth(ctx) // want `errgroup callback should probably not reference outer context "ctx", use the errgroup-derived context "egCtx"`
})
```

A *lot* more cases are covered in the [`examples.go`](testdata/base/examples.go) file!

