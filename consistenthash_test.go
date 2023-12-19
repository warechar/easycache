package easyCache

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap_Get(t *testing.T) {
	m := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	m.Add("1", "2", "3")

	fmt.Println(m.Get("3"))
}
