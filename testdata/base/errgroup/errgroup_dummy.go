package errgroup

import "context"

type Group struct{}

func (*Group) Go(func() error) {}

func (*Group) Wait() error { return nil }

func WithContext(ctx context.Context) (*Group, context.Context) {
	return new(Group), ctx
}
