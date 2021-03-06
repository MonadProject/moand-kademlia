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

func (bucket *Bucket) Empty() bool {
	return bucket.list.Len() == 0
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

// 迁移的原则是：将与本地节点更近(即​更大)节点迁移至新建 Bucket ​，
// 迁移完成后再判断新建 Bucket ​内节点数是否超过​限制，如果是，继续对该新建 Bucket ​进行分裂。
func (bucket *Bucket) Split(cpl int, local DhtID) *Bucket {
	bucket.rwl.Lock()
	defer bucket.rwl.Unlock()

	nextBucket := NewBucket()
	element := bucket.list.Front()

	for element != nil {
		id := NewDhtID(element.Value.(PeerID))
		ds := CPL(local, id)
		if ds > cpl {
			current := element
			nextBucket.list.PushFront(element.Value)
			element = element.Next()
			bucket.list.Remove(current)
		} else {
			element = element.Next()
		}
	}
	return nextBucket
}

func (bucket *Bucket) Active(id PeerID) {
	bucket.rwl.Lock()
	defer bucket.rwl.Unlock()
	element := bucket.search(id)
	if element == nil {
		return
	}

	bucket.list.MoveToFront(element)
}

//current only for method active and exist usage
func (bucket *Bucket) search(id PeerID) *list.Element {
	for element := bucket.list.Front(); element != nil; element = element.Next() {
		if element.Value.(PeerID) == id {
			return element
		}
	}
	return nil
}

func (bucket *Bucket) Exist(id PeerID) bool {
	bucket.rwl.Lock()
	defer bucket.rwl.Unlock()
	element := bucket.search(id)
	return element != nil
}

func (bucket *Bucket) Length() int {
	bucket.rwl.RLock()
	defer bucket.rwl.Unlock()
	return bucket.list.Len()
}

func (bucket *Bucket) Pop() {
	bucket.rwl.Lock()
	defer bucket.rwl.Unlock()
	last := bucket.list.Back()
	if last != nil {
		bucket.list.Remove(last)
	}
}
