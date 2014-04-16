package id

import "testing"

func TestId(t *testing.T) {
	t.Logf("This is goroutine %v", Id())
}
