package monad_kademlia

import (
	"container/list"
	"sync"
)

type Bucket struct {
	rwl  sync.RWMutex
	list *list.List
}

func NewBucket() *Bucket {
	bucket := new(Bucket)
	bucket.list = list.New()
	return bucket
}

//push front
func (bucket *Bucket) Push(id PeerID) {
	bucket.rwl.Lock()
	defer bucket.rwl.Unlock()
	bucket.list.PushFront(id)
}

func (bucket *Bucket) Remove(id PeerID) {
	bucket.rwl.Lock()
	defer bucket.rwl.Unlock()
	for element := bucket.list.Front(); element != nil; element = element.Next() {
		if element.Value.(PeerID) == id {
			bucket.list.Remove(element)
		}
	}
}

func (bucket *Bucket) AllPeers() []PeerID {
	bucket.rwl.RLock()
	defer bucket.rwl.RUnlock()
	result := make([]PeerID, 0, bucket.list.Len())

	for element := bucket.list.Front(); element != nil; element = element.Next() {
		result = append(result, element.Value.(PeerID))
	}
	return result
}
