package lru

import "testing"

type String string

func (s String) Len() int64 {
	return int64(len(s))
}

func TestGet(t *testing.T) {
	lru := New(0)

	lru.Add("key1", String("12323"))

	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "12323" {
		t.Fatalf("cache hit key1=12323 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache hit key1=12323 failed")
	}
}

func TestCache_Remove(t *testing.T) {
	lru := New(0)

	lru.Add("key1", String("12323"))

	lru.Remove("key1")

	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("remove key1 failed")
	}
}
