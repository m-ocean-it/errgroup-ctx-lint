package pkg

import (
	"context"
	"time"

	"github.com/m-ocean-it/errgroup-ctx-lint/testdata/base/errgroup"
	erGr "github.com/m-ocean-it/errgroup-ctx-lint/testdata/base/errgroup"
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

func Correct_AssignStmt_funcRunner() error {
	ctx := context.Background()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	fr := &funcRunner{}

	fr.run(func() error {
		return doSmth(ctx)
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
		return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_AssignStmt_AliasedImport() error {
	ctx := context.Background()

	eg, egCtx := erGr.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
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
		return doSmth(ctx) //nolint:abc,xyz // // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
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
		return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})

	eg.Go(func() error {
		return doSmth(egCtx)
	})

	return eg.Wait()
}

func Incorrect_DeclStmt_AliasedImport() error {
	ctx := context.Background()

	var eg, egCtx = erGr.WithContext(ctx)

	eg.Go(func() error {
		return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
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
			return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "innerEGContext"`
		})

		innerEG.Go(func() error {
			if err := doSmth(egCtx); err != nil { // want `errgroup callback must not reference outer context "egCtx", use the errgroup-derived context "innerEGContext"`
				return err
			}

			sd := smthDoer{}

			return sd.doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "innerEGContext"`
		})

		innerEG.Go(func() error {
			return doSmth(innerEGContext)
		})

		return innerEG.Wait()
	})

	eg.Go(func() error {
		return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})

	return eg.Wait()
}

func NestedErrGroup_AliasedImport() error {
	ctx := context.Background()

	eg, egCtx := erGr.WithContext(ctx)

	eg.Go(func() error {
		innerEG, innerEGContext := erGr.WithContext(egCtx)

		innerEG.Go(func() error {
			return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "innerEGContext"`
		})

		innerEG.Go(func() error {
			if err := doSmth(egCtx); err != nil { // want `errgroup callback must not reference outer context "egCtx", use the errgroup-derived context "innerEGContext"`
				return err
			}

			sd := smthDoer{}

			return sd.doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "innerEGContext"`
		})

		innerEG.Go(func() error {
			return doSmth(innerEGContext)
		})

		return innerEG.Wait()
	})

	eg.Go(func() error {
		return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})

	return eg.Wait()
}

func NoErrGroupContext() error {
	ctx := context.Background()

	eg := errgroup.New()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	eg.Go(func() error {
		return doSmth(ctx)
	})

	eg.Go(func() error {
		return doSmth(ctxWithCancel)
	})

	return eg.Wait()
}

// TryGo support

func TryGo_BasicViolation() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.TryGo(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

func TryGo_MultipleViolations() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.TryGo(func() error {
		<-ctx.Done()     // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

func MixedGoAndTryGo() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.TryGo(func() error {
		return ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

// Receiver-based context references

func CtxDone() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

func CtxErr() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		return ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

func CtxDeadline() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		_, _ = ctx.Deadline() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

func CtxValue() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		_ = ctx.Value("key") // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Non-call context references

func ChannelSend() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	ch := make(chan context.Context, 1)
	eg.Go(func() error {
		ch <- ctx // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
	<-ch
}

func SliceLiteral() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		ctxs := []context.Context{ctx} // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		_ = ctxs
		return nil
	})
	eg.Wait()
}

func MapLiteral() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		m := map[string]context.Context{"key": ctx} // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		_ = m
		return nil
	})
	eg.Wait()
}

func AssignToAny() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		var v any = ctx // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		_ = v
		return nil
	})
	eg.Wait()
}

// Multiple violations and callbacks

func MultipleViolationsOneLine() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done()     // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

func MultipleGoCallbacks() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Go(func() error {
		return ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

// Nested function literals (goroutine, defer, anonymous)

func NestedGoroutine() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		go func() {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		}()
		return nil
	})
	eg.Wait()
}

func NestedDefer() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		defer func() {
			_ = ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		}()
		return nil
	})
	eg.Wait()
}

func NestedAnonFunc() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		fn := func() {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		}
		fn()
		return nil
	})
	eg.Wait()
}

// Select pattern

func SelectPattern() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	t := time.NewTicker(time.Second)
	defer t.Stop()
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done(): // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
				return ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			case <-t.C:
			}
		}
	})
	eg.Wait()
}

// Context derivation with wrong parent

func ContextWithTimeout() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		childCtx, cancel := context.WithTimeout(ctx, time.Second) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		defer cancel()
		<-childCtx.Done()
		return nil
	})
	eg.Wait()
}

func ContextWithValue() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		valCtx := context.WithValue(ctx, "key", "val") // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		_ = valCtx
		return nil
	})
	eg.Wait()
}

// Method call with outer context

func MethodCall() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	c := &client{}
	eg.Go(func() error {
		return c.DoSomething(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

// Multiple outer contexts

func MultipleOuterContexts() {
	ctx1 := context.Background()
	ctx2 := context.TODO()
	eg, egCtx := errgroup.WithContext(ctx1)
	_ = egCtx
	eg.Go(func() error {
		<-ctx1.Done() // want `errgroup callback must not reference outer context "ctx1", use the errgroup-derived context "egCtx"`
		<-ctx2.Done() // want `errgroup callback must not reference outer context "ctx2", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Function parameter context

func FuncParamCtx(outerCtx context.Context) {
	eg, egCtx := errgroup.WithContext(outerCtx)
	_ = egCtx
	eg.Go(func() error {
		<-outerCtx.Done() // want `errgroup callback must not reference outer context "outerCtx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Context alias

func CtxAlias() {
	ctx := context.Background()
	parentCtx := ctx
	eg, egCtx := errgroup.WithContext(parentCtx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Mixed correct/incorrect

func MixedUsage() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-egCtx.Done()
		return ctx.Err() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

// Two errgroups — cross-context

func TwoErrgroups_WrongCtx() {
	ctx := context.Background()
	eg1, egCtx1 := errgroup.WithContext(ctx)
	eg2, egCtx2 := errgroup.WithContext(egCtx1)
	_ = egCtx2

	eg1.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx1"`
		return nil
	})

	eg2.Go(func() error {
		<-egCtx1.Done() // want `errgroup callback must not reference outer context "egCtx1", use the errgroup-derived context "egCtx2"`
		return nil
	})

	eg1.Wait()
	eg2.Wait()
}

func CrossErrgroup() {
	ctx := context.Background()
	eg1, egCtx1 := errgroup.WithContext(ctx)
	eg2, egCtx2 := errgroup.WithContext(ctx)
	eg1.Go(func() error {
		<-egCtx2.Done() // want `errgroup callback must not reference outer context "egCtx2", use the errgroup-derived context "egCtx1"`
		return nil
	})
	eg2.Go(func() error {
		<-egCtx1.Done() // want `errgroup callback must not reference outer context "egCtx1", use the errgroup-derived context "egCtx2"`
		return nil
	})
	eg1.Wait()
	eg2.Wait()
}

// Reassigned errgroup variable

func ReassignedErrgroup() {
	ctx := context.Background()
	eg, egCtx1 := errgroup.WithContext(ctx)
	eg, egCtx2 := errgroup.WithContext(egCtx1)
	_ = egCtx2
	eg.Go(func() error {
		<-egCtx1.Done() // want `errgroup callback must not reference outer context "egCtx1", use the errgroup-derived context "egCtx2"`
		return nil
	})
	eg.Wait()
}

// Go in loop

func GoInLoop() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	for i := 0; i < 5; i++ {
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
	}
	eg.Wait()
}

// Go conditional

func GoConditional(cond bool) {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	if cond {
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
	}
	eg.Wait()
}

// Deeply nested inside callback

func DeeplyNested() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		for i := 0; i < 3; i++ {
			if i > 0 {
				switch {
				default:
					<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
					return nil
				}
			}
		}
		return nil
	})
	eg.Wait()
}

// WithContext in nested scope (if block)

func WithContextInIf(cond bool) {
	ctx := context.Background()
	if cond {
		eg, egCtx := errgroup.WithContext(ctx)
		_ = egCtx
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
		eg.Wait()
	}
}

// Context reassigned before closure

func ReassignedCtx() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Context in panic condition

func CtxInPanic() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		if ctx.Err() != nil { // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			panic("context done")
		}
		return nil
	})
	eg.Wait()
}

// context.AfterFunc with outer ctx

func AfterFunc() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		context.AfterFunc(ctx, func() {}) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Type switch with outer ctx

func TypeSwitch() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		var v any = ctx // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		switch v.(type) {
		case context.Context:
		}
		return nil
	})
	eg.Wait()
}

// Passed to function

func PassedToFunc(f func(context.Context)) {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		f(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Other outer context (not the parent of WithContext)

func OtherOuterContext() {
	outerCtx := context.Background()
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-outerCtx.Done() // want `errgroup callback must not reference outer context "outerCtx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// TryGo nested goroutine

func TryGo_NestedGoroutine() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.TryGo(func() error {
		go func() {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		}()
		return nil
	})
	eg.Wait()
}

// Nolint on new detection patterns

func Nolint_CtxDone() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() //nolint
		return nil
	})
	eg.Wait()
}

func Nolint_ForOtherLinters_CtxDone() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() //nolint:abc,xyz // // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// All these tests should not report

func Neg_Shadow() {
	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	eg.Wait()
}

func Neg_LocalDerived() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		childCtx, cancel := context.WithCancel(egCtx)
		defer cancel()
		<-childCtx.Done()
		return childCtx.Err()
	})
	eg.Wait()
}

func Neg_LocalTimeout() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		childCtx, cancel := context.WithTimeout(egCtx, 5*time.Second)
		defer cancel()
		<-childCtx.Done()
		return childCtx.Err()
	})
	eg.Wait()
}

func Neg_LocalBackground() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		freshCtx := context.Background()
		_ = freshCtx
		return nil
	})
	eg.Wait()
}

func Neg_LocalTODO() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		todoCtx := context.TODO()
		_ = todoCtx
		return nil
	})
	eg.Wait()
}

func Neg_DiscardedCtx() {
	ctx := context.Background()
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-ctx.Done()
		return nil
	})
	eg.Wait()
}

func Neg_ParamGroup(ctx context.Context, eg *errgroup.Group) {
	eg.Go(func() error {
		<-ctx.Done()
		return nil
	})
}

type notErrgroup struct{}

func (n *notErrgroup) Go(f func() error) {}

func Neg_NonErrgroup() {
	ctx := context.Background()
	n := &notErrgroup{}
	n.Go(func() error {
		<-ctx.Done()
		return nil
	})
}

type notErrgroupTryGo struct{}

func (n *notErrgroupTryGo) TryGo(f func() error) bool { return true }

func Neg_NonErrgroupTryGo() {
	ctx := context.Background()
	n := &notErrgroupTryGo{}
	n.TryGo(func() error {
		<-ctx.Done()
		return nil
	})
}

func Neg_TwoErrgroupsCorrect() {
	ctx := context.Background()
	eg1, egCtx1 := errgroup.WithContext(ctx)
	eg2, egCtx2 := errgroup.WithContext(egCtx1)

	eg1.Go(func() error {
		<-egCtx1.Done()
		return nil
	})
	eg2.Go(func() error {
		<-egCtx2.Done()
		return nil
	})
	eg1.Wait()
	eg2.Wait()
}

func Neg_NoGo() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Wait()
}

func Neg_MultipleGoCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-egCtx.Done()
		return nil
	})
	eg.Go(func() error {
		return egCtx.Err()
	})
	eg.Go(func() error {
		childCtx, cancel := context.WithCancel(egCtx)
		defer cancel()
		<-childCtx.Done()
		return nil
	})
	eg.Wait()
}

func Neg_NestedGoroutineCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		go func() {
			<-egCtx.Done()
		}()
		return nil
	})
	eg.Wait()
}

func Neg_SelectCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	t := time.NewTicker(time.Second)
	defer t.Stop()
	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-t.C:
			}
		}
	})
	eg.Wait()
}

func Neg_MethodCallback() {
	ctx := context.Background()
	eg, _ := errgroup.WithContext(ctx)
	h := &helper{}
	eg.Go(h.work)
	eg.Wait()
}

func Neg_NamedFuncCallback() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(namedCallback)
	eg.Wait()
}

func Neg_VarFuncCallback() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	fn := func() error { return nil }
	eg.Go(fn)
	eg.Wait()
}

func Neg_ErrgroupInClosure() {
	fn := func() {
		ctx := context.Background()
		eg, egCtx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			<-egCtx.Done()
			return nil
		})
		eg.Wait()
	}
	fn()
}

type fakeCtx struct{}

func Neg_NonContextNamedCtx() {
	ctx := fakeCtx{}
	_ = ctx
	eg, egCtx := errgroup.WithContext(context.Background())
	_ = egCtx
	eg.Go(func() error {
		_ = ctx
		return nil
	})
	eg.Wait()
}

func Neg_WithSetLimit() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(5)
	eg.Go(func() error {
		<-egCtx.Done()
		return nil
	})
	eg.Wait()
}

func Neg_LocalWithDeadline() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		childCtx, cancel := context.WithDeadline(egCtx, time.Now().Add(time.Hour))
		defer cancel()
		<-childCtx.Done()
		return nil
	})
	eg.Wait()
}

func Neg_LocalWithValue() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		valCtx := context.WithValue(egCtx, "key", "val")
		_ = valCtx
		return nil
	})
	eg.Wait()
}

func Neg_EmptyCallback() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		return nil
	})
	eg.Wait()
}

func Neg_NoContextCallback() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		x := 42
		_ = x
		return nil
	})
	eg.Wait()
}

func Neg_GoInLoopCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	for i := 0; i < 5; i++ {
		eg.Go(func() error {
			<-egCtx.Done()
			return nil
		})
	}
	eg.Wait()
}

func Neg_DeeplyNestedCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		for i := 0; i < 3; i++ {
			if i > 0 {
				switch {
				default:
					<-egCtx.Done()
					return nil
				}
			}
		}
		return nil
	})
	eg.Wait()
}

func Neg_CtxInDeferCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer func() {
			_ = egCtx.Err()
		}()
		return nil
	})
	eg.Wait()
}

func Neg_MixedGoTryGoCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-egCtx.Done()
		return nil
	})
	eg.TryGo(func() error {
		return egCtx.Err()
	})
	eg.Wait()
}

func Neg_TryGoLocalDerived() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.TryGo(func() error {
		childCtx, cancel := context.WithCancel(egCtx)
		defer cancel()
		<-childCtx.Done()
		return nil
	})
	eg.Wait()
}

func Neg_MethodCallCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	c := &client{}
	eg.Go(func() error {
		return c.DoSomething(egCtx)
	})
	eg.Wait()
}

func Neg_ChannelSendCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	ch := make(chan context.Context, 1)
	eg.Go(func() error {
		ch <- egCtx
		return nil
	})
	eg.Wait()
	<-ch
}

func Neg_SliceLiteralCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		ctxs := []context.Context{egCtx}
		_ = ctxs
		return nil
	})
	eg.Wait()
}

func Neg_MapLiteralCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		m := map[string]context.Context{"key": egCtx}
		_ = m
		return nil
	})
	eg.Wait()
}

func Neg_ErrgroupInGoroutine() {
	go func() {
		ctx := context.Background()
		eg, egCtx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			<-egCtx.Done()
			return nil
		})
		eg.Wait()
	}()
}

func Neg_SequentialErrgroups() {
	ctx := context.Background()

	eg1, egCtx := errgroup.WithContext(ctx)
	eg1.Go(func() error {
		<-egCtx.Done()
		return nil
	})
	eg1.Wait()

	eg2, egCtx := errgroup.WithContext(ctx)
	eg2.Go(func() error {
		<-egCtx.Done()
		return nil
	})
	eg2.Wait()
}

func Neg_AfterFuncCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		context.AfterFunc(egCtx, func() {})
		return nil
	})
	eg.Wait()
}

func Neg_WithoutCancel() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		noCancelCtx := context.WithoutCancel(egCtx)
		_ = noCancelCtx
		return nil
	})
	eg.Wait()
}

func Neg_TryGoParamGroup(ctx context.Context, eg *errgroup.Group) {
	eg.TryGo(func() error {
		<-ctx.Done()
		return nil
	})
}

// Context used as a struct field value (unkeyed literal).
func StructLiteral() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		s := ctxHolder{ctx} // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		_ = s
		return nil
	})
	eg.Wait()
}

// Context used as a map key.
func MapKey() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		m := map[context.Context]string{ctx: "v"} // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		_ = m
		return nil
	})
	eg.Wait()
}

// Context returned from the callback.
func ReturnCtx() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		_ = ctx // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// go statement inside callback passing outer context.
func GoStmtInsideCallback() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		go doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Multiple wrong contexts in a single call.
func MultipleCtxArgs() {
	ctx1 := context.Background()
	ctx2 := context.TODO()
	eg, egCtx := errgroup.WithContext(ctx1)
	_ = egCtx
	eg.Go(func() error {
		return doSmth2(ctx1, ctx2) // want `errgroup callback must not reference outer context "ctx1", use the errgroup-derived context "egCtx"` `errgroup callback must not reference outer context "ctx2", use the errgroup-derived context "egCtx"`
	})
	eg.Wait()
}

// Context in a closure returned from the callback.
func ClosureReturnedFromCallback() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		fn := func() context.Context {
			return ctx // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		}
		_ = fn
		return nil
	})
	eg.Wait()
}

// context.WithCancel with wrong parent.
func ContextWithCancel() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		childCtx, cancel := context.WithCancel(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		defer cancel()
		<-childCtx.Done()
		return nil
	})
	eg.Wait()
}

// context.WithDeadline with wrong parent.
func ContextWithDeadline() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		childCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Hour)) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		defer cancel()
		<-childCtx.Done()
		return nil
	})
	eg.Wait()
}

// errgroup in a method receiver.
func (sd *smthDoer) ErrgroupInMethod() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// errgroup.WithContext with context.TODO() as parent.
func WithContextTODO() {
	ctx := context.TODO()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// errgroup created in switch case.
func WithContextInSwitch(n int) {
	ctx := context.Background()
	switch n {
	case 1:
		eg, egCtx := errgroup.WithContext(ctx)
		_ = egCtx
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
		eg.Wait()
	}
}

// errgroup created in for loop body.
func WithContextInLoop() {
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		eg, egCtx := errgroup.WithContext(ctx)
		_ = egCtx
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
		eg.Wait()
	}
}

// Three levels of nesting.
func TripleNestedErrGroup() {
	ctx := context.Background()
	eg1, egCtx1 := errgroup.WithContext(ctx)
	eg1.Go(func() error {
		eg2, egCtx2 := errgroup.WithContext(egCtx1)
		eg2.Go(func() error {
			eg3, egCtx3 := errgroup.WithContext(egCtx2)
			_ = egCtx3
			eg3.Go(func() error {
				<-ctx.Done()    // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx3"`
				<-egCtx1.Done() // want `errgroup callback must not reference outer context "egCtx1", use the errgroup-derived context "egCtx3"`
				<-egCtx2.Done() // want `errgroup callback must not reference outer context "egCtx2", use the errgroup-derived context "egCtx3"`
				return nil
			})
			return eg3.Wait()
		})
		return eg2.Wait()
	})
	eg1.Wait()
}

// Outer context assigned to a different variable, then used.
func CtxViaIntermediateVar() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	saved := ctx
	eg.Go(func() error {
		<-saved.Done() // want `errgroup callback must not reference outer context "saved", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Outer context in a for-range value.
func CtxInForRange() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	ctxs := []context.Context{ctx}
	eg.Go(func() error {
		for _, c := range ctxs {
			_ = c
		}
		// Still flag direct outer ctx reference.
		_ = ctx // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Context in a binary expression (comparison).
func CtxComparison() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if egCtx == ctx { // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		}
		return nil
	})
	eg.Wait()
}

// Context used both correctly and incorrectly in different Go callbacks.
func MixedCallbacks() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return doSmth(egCtx) // correct
	})
	eg.Go(func() error {
		return doSmth(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
	})
	eg.Go(func() error {
		return doSmth(egCtx) // correct
	})
	eg.Wait()
}

// --- Negative ---

// Context shadowed by a non-context variable inside the closure.
func Neg_ShadowedByNonCtx() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		ctx := 42 // shadows outer ctx with an int
		_ = ctx
		return nil
	})
	eg.Wait()
}

// errgroup with New() and separate WithContext — only New()'s Go should not flag.
func Neg_NewAndWithContext() {
	ctx := context.Background()
	plainEg := errgroup.New()
	plainEg.Go(func() error {
		return doSmth(ctx) // no flag: plainEg has no WithContext
	})
	plainEg.Wait()

	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Context derived inside the closure from the errgroup context, then passed further.
func Neg_DerivedCtxChain() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		childCtx, cancel := context.WithCancel(egCtx)
		defer cancel()
		grandchildCtx := context.WithValue(childCtx, "k", "v")
		return doSmth(grandchildCtx)
	})
	eg.Wait()
}

// errgroup in a goroutine started inside the function.
func Neg_ErrgroupInNewGoroutine() {
	ctx := context.Background()
	go func() {
		eg, egCtx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			<-egCtx.Done()
			return nil
		})
		eg.Wait()
	}()
}

// Package-level context used inside errgroup callback — should flag.
// (package-level vars are still outer scope)
var pkgCtx = context.Background()

func PkgLevelCtx() {
	eg, egCtx := errgroup.WithContext(pkgCtx)
	_ = egCtx
	eg.Go(func() error {
		<-pkgCtx.Done() // want `errgroup callback must not reference outer context "pkgCtx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// errgroup.WithContext in defer — should still track.
func WithContextInDefer() {
	ctx := context.Background()
	defer func() {
		eg, egCtx := errgroup.WithContext(ctx)
		_ = egCtx
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
		eg.Wait()
	}()
}

// errgroup.Go with multiple args (should not panic or flag — non-standard signature).
func Neg_MultipleArgsGo() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	// errgroup.Go only takes one arg; if called with zero, just don't crash.
	_ = eg
}

// Context declared as a function-level var (not short decl).
func VarDeclCtx() {
	var ctx context.Context = context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// Two errgroups in separate scopes — each should be independent.
func TwoErrgroupsSeparateScopes() {
	ctx := context.Background()
	{
		eg, egCtx := errgroup.WithContext(ctx)
		_ = egCtx
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
		eg.Wait()
	}
	{
		eg, egCtx := errgroup.WithContext(ctx)
		_ = egCtx
		eg.Go(func() error {
			<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
			return nil
		})
		eg.Wait()
	}
}

// Context used in variadic function args.
func VariadicCtxArg() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.Go(func() error {
		doSmthVariadic(ctx) // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

// TryGo with correct usage.
func Neg_TryGoCorrect() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	eg.TryGo(func() error {
		<-egCtx.Done()
		return egCtx.Err()
	})
	eg.Wait()
}

// SetLimit followed by TryGo with violation.
func TryGoWithSetLimit() {
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	_ = egCtx
	eg.SetLimit(3)
	eg.TryGo(func() error {
		<-ctx.Done() // want `errgroup callback must not reference outer context "ctx", use the errgroup-derived context "egCtx"`
		return nil
	})
	eg.Wait()
}

func doSmth(_ context.Context) error { return nil }

type smthDoer struct{}

func (sd *smthDoer) doSmth(_ context.Context) error { return nil }

type funcRunner struct{}

func (f *funcRunner) run(func() error) {}

type client struct{}

func (c *client) DoSomething(ctx context.Context) error { return nil }

type helper struct{}

func (h *helper) work() error { return nil }

func namedCallback() error { return nil }

func doSmth2(_ context.Context, _ context.Context) error { return nil }

func doSmthVariadic(_ ...context.Context) {}

type ctxHolder struct {
	ctx context.Context
}
