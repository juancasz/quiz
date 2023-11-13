package timer

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Timer interface {
	Wait(ctx context.Context)
	Finish(ctx context.Context)
}

type timer struct {
	*time.Ticker
	done chan struct{}
}

func New(seconds int) *timer {
	if seconds == 0 {
		seconds = 30
	}
	return &timer{
		Ticker: time.NewTicker(time.Duration(seconds) * time.Second),
		done:   make(chan struct{}, 1),
	}
}

func (t *timer) Wait(ctx context.Context) {
	select {
	case <-t.done:
		return
	case <-t.Ticker.C:
		fmt.Printf("\n\ntime completed\n\n")
		os.Exit(0)
	}
}

func (t *timer) Finish(ctx context.Context) {
	t.done <- struct{}{}
}
