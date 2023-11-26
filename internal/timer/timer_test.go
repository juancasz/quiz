package timer_test

import (
	"context"
	"quiz/internal/timer"
	"testing"
	"time"
)

func TestTick(t *testing.T) {
	timerTest := timer.New(1)
	ctx := context.Background()
	t.Errorf("error test tick ----------------")
	// Wait for a tick
	select {
	case <-timerTest.Tick(ctx):
		// If a tick is received, test passed
		t.Log("Received tick")
	case <-time.After(2 * time.Second):
		// If no tick or done is received after 2 seconds, test failed
		t.Error("Did not receive tick or done after 2 seconds")
	}
}
