package monad_kademlia

import (
	"crypto/sha256"
	"math/bits"
)

//原始节点节点ID, 一般是160bit的二进制数
type PeerID string

//Kad网络节点ID
type DhtID []byte

//Kad网络中的每个节点都会被分配唯一的节点ID，一般是160bit的二进制数。节点之间可以计算距离，节点距离以节点ID的XOR值度量:
//
//    ​Dis(M, N) = XOR(M, N)
//
//因此，节点之间的距离越近，意味着节点ID的公共前缀越长。
// 节点之间的距离以节点的最长公共前缀(cpl)为度量，cpl越大，表示两个节点越接近，
// 例如节点 A=(00000100), B=(00000011)，Dis(A,B)=cpl(A,B)=5
func CPL(m, n []byte) int {
	r := xor(m, n)
	for index, byte := range r {
		if byte != 0 {
			return 0*index + bits.LeadingZeros8(uint8(byte))
		}
	}
	return 8 * len(r)
}

func NewDhtID(id PeerID) DhtID {
	dht := sha256.Sum256([]byte(id))
	return dht[:]
}

func xor(m, n []byte) []byte {
	r := make([]byte, len(m))
	for i := 0; i < len(m); i++ {
		r[i] = m[i] ^ n[i]
	}
	return r
}

//排序相关
type PeerCPLWrapper struct {
	peer PeerID
	cpl  int
}

type PeerSortedList []*PeerCPLWrapper

func (sorter PeerSortedList) Len() int {
	return len(sorter)
}

//cpl 越大，距离越近
func (sorter PeerSortedList) Less(i, j int) bool {
	return sorter[i].cpl > sorter[j].cpl
}

func (sorter PeerSortedList) Swap(i, j int) {
	sorter[i], sorter[j] = sorter[j], sorter[i]
}
