package util

import (
	"testing"
)

func TestCircuitBreaker1(t *testing.T) {
	cb := NewCircuitBreaker()
	t.Log(cb.Allow())
}
