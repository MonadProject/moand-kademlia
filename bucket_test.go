package monad_kademlia

import (
	"fmt"
	"testing"
)

func TestNewBucket(t *testing.T) {
	bucket := NewBucket()
	fmt.Printf("bucket is %v",bucket)
}

