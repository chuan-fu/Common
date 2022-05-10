// Copyright (c) 2016,2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ratelimit // import "go.uber.org/ratelimit"

import (
	"context"
	"time"

	"github.com/andres-erbsen/clock"
)

// Note: This file is inspired by:
// https://github.com/prashantv/go-bench/blob/master/ratelimit

// Limiter is used to rate-limit some process, possibly across goroutines.
// The process is expected to call Take() before every iteration, which
// may block to throttle the goroutine.
type Limiter interface {
	// Take should block to make sure that the RPS is met.
	Take() time.Time
	TakeWithContext(ctx context.Context) time.Time
	TakeWithTimeOut(timeout time.Duration) time.Time
	TakeWithOpt(opt TimeOutOption) time.Time
}

// Clock is the minimum necessary interface to instantiate a rate limiter with
// a clock or mock clock, compatible with clocks created using
// github.com/andres-erbsen/clock.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

// config configures a limiter.
type config struct {
	clock    Clock
	slack    int
	per      time.Duration
	maxSleep time.Duration // 最大睡眠时间
}

// New returns a Limiter that will limit to the given RPS.
func New(rate int, opts ...Option) Limiter {
	return newAtomicBased(rate, opts...)
}

// buildConfig combines defaults with options.
func buildConfig(opts []Option) config {
	c := &config{
		clock:    clock.New(),
		slack:    10,
		per:      time.Second,
		maxSleep: 0, // 默认无超时
	}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}

// Option configures a Limiter.
type Option func(c *config)

// WithClock returns an option for ratelimit.New that provides an alternate
// Clock implementation, typically a mock Clock for testing.
func WithClock(clock Clock) Option {
	return func(c *config) {
		c.clock = clock
	}
}

// WithoutSlack configures the limiter to be strict and not to accumulate
// previously "unspent" requests for future bursts of traffic.
var WithoutSlack = WithSlack(0)

// WithSlack configures custom slack.
// Slack allows the limiter to accumulate "unspent" requests
// for future bursts of traffic.
func WithSlack(slack int) Option {
	return func(c *config) {
		c.slack = slack
	}
}

type perOption time.Duration

func (p perOption) apply(c *config) {
	c.per = time.Duration(p)
}

// WithPer allows configuring limits for different time windows.
//
// The default window is one second, so New(100) produces a one hundred per
// second (100 Hz) rate limiter.
//
// New(2, WithPer(60*time.Second)) creates a 2 per minute rate limiter.
func WithPer(per time.Duration) Option {
	return func(c *config) {
		c.per = per
	}
}

func WithMaxSleep(maxSleep time.Duration) Option {
	return func(c *config) {
		c.maxSleep = maxSleep
	}
}

type unlimited struct{}

// NewUnlimited returns a RateLimiter that is not limited.
func NewUnlimited() Limiter {
	return unlimited{}
}

func (unlimited) Take() time.Time {
	return time.Now()
}

func (unlimited) TakeWithContext(ctx context.Context) time.Time {
	return time.Now()
}

func (unlimited) TakeWithTimeOut(timeout time.Duration) time.Time {
	return time.Now()
}

func (unlimited) TakeWithOpt(opt TimeOutOption) time.Time {
	return time.Now()
}
