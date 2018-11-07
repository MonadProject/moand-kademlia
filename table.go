package monad_kademlia

import (
	"sort"
	"sync"
)

//每次查询最近节点时，返回的结果数
const PeerCount = 10

//构建路由表的本质是建立到网络全局的地图，目标是：对于节点​ M ，给定任意节点 X ​，可以根据节点​很容易计算出​距离 X 更近的节点列表
// 虽然我们的目标是本地一步到位的查找，但这是不现实的，这需要维护数量巨大的全局节点信息
// 我们退而求其次，采用迭代查找的思路：每次查找要距离目标更近一点点

//路由表
type Table struct {
	self       DhtID
	Buckets    []*Bucket
	bucketsize int
	rwl        sync.RWMutex
}

//size 一般为20
func NewTable(self DhtID, size int) *Table {
	return &Table{
		self:       self,
		Buckets:    []*Bucket{NewBucket()},
		bucketsize: size,
	}
}

func (table *Table) Find(id PeerID) PeerSortedList {
	table.rwl.RLock()
	cpl := CPL(table.self, NewDhtID(id))
	if cpl >= len(table.Buckets) {
		cpl = len(table.Buckets) - 1
	}
	bucket := table.Buckets[cpl]

	list := make(PeerSortedList, 0, PeerCount)
	for element := bucket.list.Front(); element != nil; element = element.Next() {
		peerWrapper := &PeerCPLWrapper{
			peer: element.Value.(PeerID),
			cpl:  CPL(table.self, NewDhtID(element.Value.(PeerID))),
		}

		if len(list) == PeerCount {
			break
		}
		list = append(list, peerWrapper)
	}

	//todo if the length is not enough

	table.rwl.Unlock()

	sort.Sort(list)

	return list
}

func (table *Table) Add(peer PeerID) {
	table.rwl.Lock()
	defer table.rwl.Unlock()
	innerId := NewDhtID(peer)

	//compute cpl
	cpl := CPL(table.self, innerId)

	if cpl >= len(table.Buckets) {
		cpl = len(table.Buckets) - 1
	}

	bucket := table.Buckets[cpl]

	if bucket.Exist(peer) {
		bucket.Active(peer)
		return
	}

	bucket.list.PushFront(peer)

	if bucket.Length() > table.bucketsize {
		if cpl == len(table.Buckets)-1 {
			//split it！
			table.split()
		} else {
			bucket.Pop()
		}
	}
}

func (table *Table) split() {
	currentBucket := table.Buckets[len(table.Buckets)-1]
	nextBucket := currentBucket.Split(len(table.Buckets)-1, table.self)
	//if split have done nothing
	if nextBucket.Empty() {
		currentBucket.Pop()
		return
	}

	table.Buckets = append(table.Buckets, nextBucket)
	if nextBucket.Length() > table.bucketsize {
		table.split()
	}
}
