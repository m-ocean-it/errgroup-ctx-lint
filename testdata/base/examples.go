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
		return doSmth(ctx) // want "passing non-errgroup context to function within errgroup-goroutine while there is an errgroup-context defined"
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt_Nolint() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) //nolint
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt_Nolint_ErrGroupCtxLint() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) //nolint:errgroup-ctx-lint
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt_Nolint_ErrGroupCtxLint_WithOtherLinters() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) //nolint:abc,all,xyz
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt_Nolint_All() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) //nolint:all
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt_Nolint_All_WithOtherLinters() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) //nolint:abc,all,xyz
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt_Nolint_ForOtherLinters() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) //nolint:abc,xyz // // want "passing non-errgroup context to function within errgroup-goroutine while there is an errgroup-context defined"
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
		return doSmth(ctx) // want "passing non-errgroup context to function within errgroup-goroutine while there is an errgroup-context defined"
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func NestedErrGroup() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		innerEG, innerEGContext := errgroup.WithContext(egCtx)

		innerEG.Go(func() error {
			return doSmth(ctx) // want "passing non-errgroup context to function within errgroup-goroutine while there is an errgroup-context defined"
		})

		innerEG.Go(func() error {
			return doSmth(egCtx) // want "passing non-errgroup context to function within errgroup-goroutine while there is an errgroup-context defined"
		})

		innerEG.Go(func() error {
			return doSmth(innerEGContext)
		})

		return innerEG.Wait()
	})

	eg.Go(func() error {
		return doSmth(ctx) // want "passing non-errgroup context to function within errgroup-goroutine while there is an errgroup-context defined"
	})

	return eg.Wait()
}

// TODO:
// func NoErrGroupContext() error {
// 	ctx := context.Background()

// 	eg := errgroup.New()
// 	ctxWithCancel, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	eg.Go(func() error {
// 		return doSmth(ctx)
// 	})

// 	eg.Go(func() error {
// 		return doSmth(ctxWithCancel)
// 	})

// 	return eg.Wait()
// }

func doSmth(_ context.Context) error { return nil }
