package util

import (
	"sync/atomic"
)

const (
	kCircuitOpen        = 0
	kCircuitHalfOpen    = 1
	kCircuitClose       = 2
	kCalcStatIntervalMs = 1000
)

type CircuitBreaker struct {
	forceClose               bool
	forceOpen                bool
	errorThresholdPercentage float32
	sleepWindowsInMs         int64
	lastCalcStatTimeMs       int64
	lastCircuitOpenTime      int64
	successCount             int64
	failCount                int64
	circuitStatus            int8
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		forceClose:               false,
		forceOpen:                false,
		errorThresholdPercentage: 0.5,
		sleepWindowsInMs:         200,
		lastCalcStatTimeMs:       0,
		lastCircuitOpenTime:      0,
		successCount:             0,
		failCount:                0,
		circuitStatus:            kCircuitClose,
	}
}

// Allow returns true if a request is within the circuit breaker norms.
// Otherwise, it returns false.
func (cb *CircuitBreaker) Allow() bool {
	// force open the circuit, link is break so this is not allowed.
	if cb.forceOpen {
		return false
	}
	// force close the circuit, link is not break so this is allowed.
	if cb.forceClose {
		return true
	}

	var now_ms int64
	now_ms = NowInMs()
	cb.CalcStat(now_ms)

	if cb.circuitStatus == kCircuitClose {
		return true
	} else {
		if cb.IsAfterSleepWindow(now_ms) {
			cb.lastCircuitOpenTime = now_ms
			cb.circuitStatus = kCircuitHalfOpen
			// sleep so long time, try ones, and set status to half-open
			return true
		}
	}
	return false
}

// set percentage of error threshold.
func (cb *CircuitBreaker) SetErrorThresholdPercentage(percentage float32) {
	cb.errorThresholdPercentage = percentage
}

func (cb *CircuitBreaker) SetForceOpen() {
	cb.forceOpen = true
	cb.forceClose = false
}

func (cb *CircuitBreaker) SetForceClose() {
	cb.forceOpen = false
	cb.forceClose = true
}

// call me when the service called success.
func (cb *CircuitBreaker) FeedSuccess() {
	atomic.AddInt64(&cb.successCount, 1)
}

// call me when the service called failed.
func (cb *CircuitBreaker) FeedFail() {
	atomic.AddInt64(&cb.failCount, 1)
}

// check whether after the sleep windows
func (cb *CircuitBreaker) IsAfterSleepWindow(now_ms int64) bool {
	return now_ms > cb.lastCircuitOpenTime+cb.sleepWindowsInMs
}

func (cb *CircuitBreaker) CalcStat(now_ms int64) {
	if now_ms > cb.lastCalcStatTimeMs+kCalcStatIntervalMs {
		allCount := cb.successCount + cb.failCount
		if allCount > 0 {
			rate := float32(cb.failCount) / float32(allCount)
			if cb.failCount > 0 && rate >= cb.errorThresholdPercentage {
				// mark CLOSE to OPEN
				cb.circuitStatus = kCircuitOpen
				cb.lastCircuitOpenTime = now_ms
			} else {
				// mark OPEN to CLOSE
				cb.lastCircuitOpenTime = 0
				cb.circuitStatus = kCircuitClose
			}
		}

		// clear count
		atomic.StoreInt64(&cb.successCount, 0)
		atomic.StoreInt64(&cb.failCount, 0)
		cb.lastCalcStatTimeMs = now_ms
	}
}
