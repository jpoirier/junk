package throw

import "testing"

func TestThrow(t *testing.T) {
	defer func() {
		recover()
	}()
	Throw()
}
