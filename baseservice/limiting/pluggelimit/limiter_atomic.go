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

package pluggelimit // import "go.uber.org/ratelimit"

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"
)

type state struct {
	last     time.Time
	sleepFor time.Duration
}

type atomicLimiter struct {
	state unsafe.Pointer
	//lint:ignore U1000 Padding is unused but it is crucial to maintain performance
	// of this rate limiter in case of collocation with other frequently accessed memory.
	padding [56]byte // cache line size - state pointer size = 64 - 8; created to avoid false sharing.

	perRequest time.Duration
	maxSlack   time.Duration
	timeoutOpt TimeOutOption // 最大睡眠时间，<=0则无超时
	clock      Clock
}

// newAtomicBased returns a new atomic based limiter.
func newAtomicBased(rate int, opts ...Option) *atomicLimiter {
	// TODO consider moving config building to the implementation
	// independent code.
	cfg := buildConfig(opts)
	perRequest := cfg.per / time.Duration(rate)
	l := &atomicLimiter{
		perRequest: perRequest,
		maxSlack:   -1 * time.Duration(cfg.slack) * perRequest,
		timeoutOpt: WithDurationTimeOutOption(cfg.maxSleep),
		clock:      cfg.clock,
	}

	initialState := state{
		last:     time.Time{},
		sleepFor: 0,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&initialState))
	return l
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *atomicLimiter) Take() time.Time {
	return t.TakeWithOpt(t.timeoutOpt)
}

func (t *atomicLimiter) TakeWithContextDeadline(ctx context.Context) time.Time {
	return t.TakeWithOpt(WithCtxTimeOutOption(ctx))
}

func (t *atomicLimiter) TakeWithTimeOut(timeout time.Duration) time.Time {
	return t.TakeWithOpt(WithDurationTimeOutOption(timeout))
}

func (t *atomicLimiter) TakeWithOpt(opt TimeOutOption) time.Time {
	var (
		newState state
		taken    bool
		interval time.Duration
	)
	for !taken {
		now := t.clock.Now()

		previousStatePointer := atomic.LoadPointer(&t.state)
		oldState := (*state)(previousStatePointer)

		newState = state{
			last:     now,
			sleepFor: oldState.sleepFor, // maxSlack <= newState.sleepFor <= 0
		}

		// If this is our first request, then we allow it.
		if oldState.last.IsZero() {
			taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
			continue
		}

		// sleepFor calculates how much time we should sleep based on
		// the perRequest budget and how long the last request took.
		// Since the request may take longer than the budget, this number
		// can get negative, and is summed across requests.

		// 最后执行时间-当前时间+时间间隔，就是需要睡眠时间/下次执行时间
		newState.sleepFor += t.perRequest - now.Sub(oldState.last)
		// We shouldn't allow sleepFor to get too negative, since it would mean that
		// a service that slowed down a lot for a short period of time would get
		// a much higher RPS following that.
		// 睡眠时间如果小于maxSlack，则设为maxSlack，默认maxSlack为负数
		if newState.sleepFor < t.maxSlack {
			newState.sleepFor = t.maxSlack
		}
		if newState.sleepFor > 0 {
			// 睡眠时间大于0，就把睡眠时间+当前时间，代表最后执行时间
			newState.last = newState.last.Add(newState.sleepFor)
			interval, newState.sleepFor = newState.sleepFor, 0
		}

		// 当maxSlack为负数时，可能存在 maxSlack <= newState.sleepFor <= 0
		// 此时会将newState.sleepFor写入全局，睡眠时间interval为0，代表直接返回
		// 并会将newState.sleepFor写入下个Take()的睡眠时间
		// 例如：rate为10，执行间隔为100ms，上次执行时间默认为maxSlack：-100ms，则当前睡眠时间为0，直接返回

		// 如果设置了最大睡眠时间，且需要睡眠大于最大睡眠时间，直接返回
		// 可能存在多个协程竞争t.state的写入，竞争更少的写入时间，如果竞争失败，则需要更长的写入时间
		// 如果当前interval都超过maxSleep，失败的话，需要的时间更长，所以不需要等到写入成功在返回失败
		// 而且如果写入成功再返回超时，当前协程也被算进时间里了，但是当前并不执行，会造成大量失败
		if opt != nil && opt(&newState, interval) {
			return time.Time{}
		}

		// taken == false，写入失败，代表锁竞争，重新循环尝试写入
		taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}

	t.clock.Sleep(interval)
	return newState.last
}

type TimeOutOption func(s *state, interval time.Duration) bool // true表示已过期

func WithCtxTimeOutOption(ctx context.Context) TimeOutOption {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil
	}
	return func(s *state, interval time.Duration) bool {
		return s.last.After(deadline)
	}
}

func WithDurationTimeOutOption(timeout time.Duration) TimeOutOption {
	if timeout <= 0 {
		return nil
	}
	return func(s *state, interval time.Duration) bool {
		return interval > timeout
	}
}
