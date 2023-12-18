package easyCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpPool struct {
	self     string
	basePath string
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
