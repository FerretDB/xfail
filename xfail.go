// Package xfail provides testing helpers for expected tests failures.
package xfail

import (
	"fmt"
	"math"
	"sync/atomic"
	"testing"
)

// XFail returns a new [testing.TB] instance that expects the test to fail for the given reason.
//
// At the end of the test, if it was marked as failed using non-fatal methods Fail, Error, or Errorf,
// it will pass instead.
// If it wasn't marked as failed, it will fail so that XFail call can be removed.
// Fatal methods FailNow, Fatal, or Fatalf will skip the rest of the test instead.
//
// It panics if the amount of marked tests is larger than math.MaxInt64 (9223372036854775807).
func XFail(tb testing.TB, reason string) testing.TB {
	tb.Helper()

	if reason == "" {
		tb.Fatal("XFail reason can't be empty")
	}

	x := &xfail{
		TB: tb,
	}

	x.TB.Cleanup(func() {
		if x.failed.Load() {
			x.TB.Logf("Test failed as expected: %s", reason)

			if count := failed.Add(1); count < 0 {
				panic(fmt.Sprintf("The number of failed tests exceeded the max int64 range (%d)", math.MaxInt64))
			}

			return
		}

		x.TB.Fatalf("Test passed unexpectedly: %s", reason)
	})

	return x
}

// xfail wraps [testing.TB] with expected failure logic.
type xfail struct {
	testing.TB
	failed atomic.Bool // we can't access testing.common.failed/skipped/etc fields
}

// Failed reports whether the function has failed.
func (x *xfail) Failed() bool {
	return x.failed.Load()
}

// Fail marks the function as having failed but continues execution.
func (x *xfail) Fail() {
	x.failed.Store(true)
}

// Error is equivalent to Log followed by Fail.
func (x *xfail) Error(args ...any) {
	x.Log(args...)
	x.Fail()
}

// Errorf is equivalent to Logf followed by Fail.
func (x *xfail) Errorf(format string, args ...any) {
	x.Logf(format, args...)
	x.Fail()
}

// FailNow marks the function as having failed and stops its execution.
func (x *xfail) FailNow() {
	x.Fail()

	// runtime.Goexit would not work
	x.SkipNow()
}

// Fatal is equivalent to Log followed by FailNow.
func (x *xfail) Fatal(args ...any) {
	x.Log(args...)
	x.FailNow()
}

// Fatalf is equivalent to Logf followed by FailNow.
func (x *xfail) Fatalf(format string, args ...any) {
	x.Logf(format, args...)
	x.FailNow()
}

// failed stores the amount of failed tests marked by XFail().
var failed atomic.Int64

// FailedCount returns the number of failed tests that were expected to fail.
//
// It instantly returns the count at the time of being called, hence it's a caller responsibility
// to make sure that all tests finished.
func FailedCount() int64 {
	return failed.Load()
}

// check interfaces
var (
	_ testing.TB = (*xfail)(nil)
)
