package ratelimit

import (
	"context"
	"time"
)

type tokens chan struct{}

type TokenBucket struct {
	Tokens tokens
	count  int64
	ticker *time.Ticker
}

func NewTokenBucket(count, rate int64) *TokenBucket {
	tokens := make(chan struct{}, count)

	c := int(count)
	for i := 0; i < c; i++ {
		tokens <- struct{}{}
	}

	everyMs := 1 / float64(rate) * 1000
	return &TokenBucket{
		Tokens: tokens,
		count:  count,
		ticker: time.NewTicker(time.Duration(int64(everyMs) * int64(time.Millisecond))),
	}
}

func (tb *TokenBucket) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-tb.ticker.C:
				select {
				case tb.Tokens <- struct{}{}:
				default:
				}
			case <-ctx.Done():
				// The context was likely cancelled
				return
			}
		}
	}()
}

func (tb *TokenBucket) Consume() {
	<-tb.Tokens
}
