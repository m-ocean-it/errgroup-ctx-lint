package pkg

import (
	"context"

	"github.com/m-ocean-it/errgroup-ctx-lint/testdata/base/errgroup"
)

func Correct_AssignStmt() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Correct_DeclStmt() error {
	ctx := context.Background()

	var eg, egCtx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) // want "passing non-errgroup context to function withing errgroup-goroutine while there is an errgroup-context defined"

	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_DeclStmt() error {
	ctx := context.Background()

	var eg, egCtx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) // want "passing non-errgroup context to function withing errgroup-goroutine while there is an errgroup-context defined"

	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func doSmth(_ context.Context) error { return nil }
