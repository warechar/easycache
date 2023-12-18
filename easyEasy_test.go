package easyCache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

func TestGroup_Get(t *testing.T) {
	g := NewGroup("score", 0, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[other storage] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}

		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
	}
}
