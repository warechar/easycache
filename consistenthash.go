package easyCache

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // hash func
	replicas int            // virtual host count
	keys     []int          // save all node , include virtual node
	hashMap  map[int]string // save true node
}

func New(replicas int, hash Hash) *Map {

}
