package thqcache

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var db = map[string]string{
	"thq1": "1",
	"thq2": "2",
	"thq3": "3",
}

func TestGroup_Get(t *testing.T) {
	loadCount := make(map[string]int, len(db))
	getter := func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			if _, ok := loadCount[key]; !ok {
				loadCount[key] = 0
			}
			loadCount[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}

	g := NewGroup("test", 2<<10, GetterFunc(getter))

	for k, v := range db {
		value, err := g.Get(k)
		require.NoError(t, err)
		require.Equal(t, value.String(), v)
		_, err = g.Get(k)
		require.NoError(t, err)
		require.True(t, loadCount[k] == 1)
	}
	_, err := g.Get("no")
	require.True(t, err != nil)
}
