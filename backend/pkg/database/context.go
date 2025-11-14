package database

import (
	"context"
	"time"
)

const (
	ShortTimeout  = 3 * time.Second
	MediumTimeout = 5 * time.Second
	LongTimeout   = 10 * time.Second
)

func WithShortTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), ShortTimeout)
}

func WithMediumTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), MediumTimeout)
}

func WithLongTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), LongTimeout)
}
