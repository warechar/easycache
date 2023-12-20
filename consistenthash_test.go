package easyCache

import (
	"strconv"
	"testing"
)

func TestMap_Get(t *testing.T) {
	m := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	// 1 2 3 11 12 13 21 22 23
	m.Add("1", "2", "3")

	testM := map[string]string{
		"1":  "1",
		"22": "2",
		"27": "1",
	}

	for k, v := range testM {
		if m.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	m.Add("8")

	testM["27"] = "8"

	for k, v := range testM {
		if m.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
