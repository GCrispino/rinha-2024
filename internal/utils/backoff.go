package utils

import (
	"time"

	backoff "github.com/cenkalti/backoff/v4"
)

func DefaultBackoff() backoff.BackOff {
	boff := backoff.NewExponentialBackOff()
	boff.MaxElapsedTime = 5 * time.Second

	return boff
}
