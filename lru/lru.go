package lru

import "container/list"

/**
based on lru implement,depend on container/list package
*/

type Cache struct {
	maxBytes int64
	useBytes int64
	ll       *list.List
	cache    map[string]*list.Element
}

// 元素体
type entry struct {
	key   string
	value Value
}

// Value interface implode how many bytes it take
type Value interface {
	Len() int64
}

// New is construct of cache
func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		useBytes: 0,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

// Get search key from cache
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		return kv.value, true
	}

	return
}

func (c *Cache) Remove(key string) {
	if ele, ok := c.cache[key]; ok {
		c.ll.Remove(ele)
		delete(c.cache, key)
		c.useBytes -= int64(len(key)) + ele.Value.(*entry).value.Len()
	}
}

// removeOldest delete cache key from cache
func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)

		kv := ele.Value.(*entry)

		delete(c.cache, kv.key) // delete cache map

		c.useBytes -= int64(len(kv.key)) + kv.value.Len()
	}
}

func (c *Cache) Add(key string, value Value) {
	// if exist key
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		c.useBytes += value.Len() - kv.value.Len()
		kv.value = value // save new value
	} else {
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.useBytes += int64(len(key)) + value.Len()
	}

	for c.isFull() {
		c.removeOldest()
	}
}

func (c *Cache) isFull() bool {
	if c.maxBytes == 0 {
		return false
	}

	if c.useBytes > c.maxBytes {
		return true
	}

	return false
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
