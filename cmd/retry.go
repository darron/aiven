package cmd

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// GetBackoff returns a backoff.BackOff that we can use to retry operations inline.
func GetBackoff(interval, maxTime time.Duration) (backoff.BackOff, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), maxTime)
	b := backoff.NewConstantBackOff(interval)
	return backoff.WithContext(b, ctx), cancel
}
