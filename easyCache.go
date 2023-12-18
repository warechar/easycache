package easyCache

import (
	"errors"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string // unique group name
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()

	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}

	groups[name] = g
	return g
}

// GetGroup returns the named group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[easyCache] hit " + key)
		return v, nil
	}

	// load value from getter if cache not found
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	clone := make([]byte, len(bytes))
	copy(clone, bytes) // protect cache data if input data updated
	value := ByteView{b: clone}
	g.mainCache.set(key, value)

	return value, nil
}
