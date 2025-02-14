// Implement a rate limiter using goroutines and channels

package main

import (
	"fmt"
	"time"
)

// RateLimiter represents a token bucket rate limiter
type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	done   chan bool
	rate   int
	burst  int
}

// NewRateLimiter creates a new rate limiter with specified requests per second and burst size
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, burst),
		done:   make(chan bool),
		rate:   rate,
		burst:  burst,
	}

	// Initialize with just one token instead of full burst
	rl.tokens <- struct{}{}

	// Start token refill goroutine with precise timing
	go rl.refillTokens()
	return rl
}

// refillTokens periodically adds tokens to the bucket
func (rl *RateLimiter) refillTokens() {
	// More precise timing for token refill
	ticker := time.NewTicker(time.Second / time.Duration(rl.rate))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
				// Token added
			default:
				// Bucket is full
			}
		case <-rl.done:
			return
		}
	}
}

// Allow checks if a request can proceed
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	rl.done <- true
}

func main() {
	// Create a rate limiter with 2 requests per second and burst size of 3
	limiter := NewRateLimiter(2, 3)
	defer limiter.Stop()

	start := time.Now()

	// Simulate requests
	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Printf("%.2fs: Request %d allowed\n", time.Since(start).Seconds(), i)
		} else {
			fmt.Printf("%.2fs: Request %d throttled\n", time.Since(start).Seconds(), i)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
