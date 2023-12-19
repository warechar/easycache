package easyCache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // hash func
	replicas int            // virtual host count
	keys     []int          // save all node , include virtual node
	hashMap  map[int]string // save true node
}

func New(replicas int, hash Hash) *Map {
	map1 := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if map1.hash == nil {
		map1.hash = crc32.ChecksumIEEE // default hash
	}

	return map1
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// look for m.keys
	//  if idx == len(m.keys), should be % len(m.keys)
	// sort.Search not found, return len(m.keys)
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx]%len(m.keys)]
}
