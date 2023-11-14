package timer

import (
	"context"
	"time"
)

type Tick <-chan time.Time
type Done chan struct{}

type Timer interface {
	Tick(ctx context.Context) Tick
	Done(ctx context.Context) Done
	Finish(ctx context.Context)
}

type timer struct {
	*time.Ticker
	done Done
}

func New(seconds int) *timer {
	return &timer{
		Ticker: time.NewTicker(time.Duration(seconds) * time.Second),
		done:   make(chan struct{}, 1),
	}
}

func (t *timer) Tick(ctx context.Context) Tick {
	return t.Ticker.C
}

func (t *timer) Done(ctx context.Context) Done {
	return t.done
}

func (t *timer) Finish(ctx context.Context) {
	t.done <- struct{}{}
}
