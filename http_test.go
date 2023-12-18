package easyCache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestServer_ServeHTTP(t *testing.T) {
	NewGroup("scores", 0, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[db] search key" + key)

		if v, ok := db[key]; ok {
			return []byte(v), nil
		}

		return nil, fmt.Errorf("%s not1111 exist", key)
	}))

	addr := "127.0.0.1:8080"
	httppool := NewHttpPool(addr)

	http.ListenAndServe(addr, httppool)
}
