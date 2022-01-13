package lru

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_Add(t *testing.T) {
	lru := NewCache(0, nil)
	key, value := "key", String("test")
	lru.Add(key, value)
	v, ok := lru.Get(key)
	require.True(t, ok)
	require.Equal(t, v, value)
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key", "k3"
	v1, v2, v3 := "value1", "value", "v3"
	maxLen := len(k1 + k2 + v1 + v2)
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := NewCache(int64(maxLen), callback)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	v, ok := lru.Get(k1)
	require.False(t, ok)
	v, ok = lru.Get(k2)
	require.True(t, ok)
	require.Equal(t, v, String(v2))
	require.Equal(t, 2, lru.Len())
	require.Equal(t, int64(len(k2+k3+v2+v3)), lru.nowBytes)
	require.Equal(t, []string{k1}, keys)
}
