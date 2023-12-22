package easyCache

import (
	"errors"
	"log"
	"sync"
)

// Getter getter data for a key
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
	peers     PeerPicker
}

var (
	mu     sync.RWMutex
	groups = sync.Map{}
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

	groups.Store(name, g)
	return g
}

// GetGroup returns the named group
func GetGroup(name string) *Group {
	v, _ := groups.Load(name)
	return v.(*Group)
}

// get local value
func (g *Group) getLocally(key string) (ByteView, error) {
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

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}

	// get locally cached first
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[easyCache] hit " + key)
		return v, nil
	}

	// if local cache does not exist, the cache is obtained from the remote peer
	return g.load(key)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("register called more than once")
	}

	g.peers = peers
}

func (g *Group) load(key string) (ByteView, error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err := g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
		}
	}

	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {

	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: bytes}, nil
}
