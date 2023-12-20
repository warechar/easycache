package easyCache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const defaultReplicas = 50

type HttpPool struct {
	self       string
	basePath   string
	peers      *Map
	mu         sync.Mutex
	httpGetter map[string]*httpGetter
}

type httpGetter struct {
	baseUrl string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseUrl, group, key)

	resp, err := http.Get(u)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

func time33(data []byte) uint32 {
	hash := int32(5381)
	m := md5.New()
	md5Str := hex.EncodeToString(m.Sum(data))

	for i := 0; i < 32; i++ {
		hash += hash<<5 + int32(md5Str[i])
	}

	return uint32(hash & 0x7FFFFFFF)
}

func (h *HttpPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.peers = New(defaultReplicas, time33)
	h.peers.Add(peers...) // add peers
	h.httpGetter = make(map[string]*httpGetter, len(peers))

	for _, peer := range peers {
		h.httpGetter[peer] = &httpGetter{baseUrl: peer}
	}
}

func (h *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		h.Log("pick peer %s", peer)
		return h.httpGetter[peer], true
	}

	return nil, false
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: "/_easycahe/",
	}
}

func (h *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", h.self, fmt.Sprintf(format, v...))
}

// ServerHTTP implement http.Handler
func (h *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log("%s %s", r.Method, r.URL.Path)
	if !strings.HasPrefix(r.URL.Path, "/_easycache/") {
		//panic("HTTPPool serving unexpected path: " + r.URL.Path)
		return
	}

	parts := strings.SplitN(r.URL.Path[len("/_easycache/"):], "/", 2)

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group :"+groupName, http.StatusBadRequest)
		return
	}

	v, err := group.Get(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(v.String()))
}
