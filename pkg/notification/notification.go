package notification

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Notifier interface {
	Notify(ctx context.Context, message string) error
}

type Notifiers []Notifier

func (ns Notifiers) Notify(ctx context.Context, message string) error {
	switch l := len(ns); l {
	case 0:
		return nil
	case 1:
		return ns[0].Notify(ctx, message)
	default:
		var group errgroup.Group
		for _, n := range ns {
			n := n
			group.Go(func() error {
				return n.Notify(ctx, message)
			})
		}
		return group.Wait()
	}
}
