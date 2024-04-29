package xfail

import (
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	// comment out to test manually
	t.Skip("test fails (as expected)")

	t.Parallel()

	t.Run("Fatal", func(t *testing.T) {
		t.Fatal("Fatal")

		panic("not reached")
	})

	if !t.Failed() {
		t.Fatal("test should be marked as failed")
	}

	t.Run("ErrorAndSkip", func(t *testing.T) {
		t.Error("Error")
		t.SkipNow()

		panic("not reached")
	})

	if !t.Failed() {
		t.Fatal("test should be marked as failed")
	}

	t.Skip("skipping failing test does not mark it as not failed")
}

func TestXFail(t *testing.T) {
	t.Parallel()

	t.Run("Fatal", func(tt *testing.T) {
		t := XFail(tt, "expected failure")
		t.Fatal("Fatal")

		panic("not reached")
	})

	if t.Failed() {
		t.Fatal("test should not be marked as failed")
	}

	t.Run("ErrorAndSkip", func(tt *testing.T) {
		t := XFail(tt, "expected failure")
		t.Error("Error")
		t.SkipNow()

		panic("not reached")
	})

	if t.Failed() {
		t.Fatal("test should not be marked as failed")
	}
}

func TestREADMEParseDuration(tt *testing.T) {
	t := XFail(tt, "https://github.com/golang/go/issues/67076")

	if _, err := time.ParseDuration("3.336e-6s"); err != nil {
		t.Fatal(err)
	}
}
