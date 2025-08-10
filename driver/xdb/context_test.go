package xdb

import (
	"context"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	timeout := 1 * time.Second
	sleepSecond := 2 * time.Second
	ctxCancel, _ := context.WithTimeout(ctx, timeout)
	t.Log("Goroutine is hanging")
	time.Sleep(sleepSecond)
	t.Log("Goroutine is awakened")
	err := Context(ctxCancel, db)
	if err != nil {
		t.Fatalf("Err(%+v)", err)
	}
}
