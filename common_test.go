package monad_kademlia

import (
	"testing"
)

func TestCPL(t *testing.T) {
	m := []byte{4}
	n := []byte{3}
	cpl := CPL(m, n)
	if cpl != 5 {
		t.Error("result is error")
	}
}
