package timer

import (
	"context"
	"time"
)

type Tick <-chan time.Time

type Timer interface {
	Tick(ctx context.Context) Tick
}

type timer struct {
	*time.Ticker
}

func New(seconds int) *timer {
	return &timer{
		Ticker: time.NewTicker(time.Duration(seconds) * time.Second),
	}
}

func (t *timer) Tick(ctx context.Context) Tick {
	return t.Ticker.C
}
