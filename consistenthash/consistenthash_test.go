package consistenthash

import (
	"github.com/stretchr/testify/require"
	"hash/crc32"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		v := crc32.ChecksumIEEE(data)
		t.Log(v)
		return v
	})

	hash.Add("http://127.0.0.1:1", "http://127.0.0.1:10", "http://127.0.0.1:100")
	//虚拟节点[771555706 1060649185 1143388674 1406271808 2186562736 2217771288 2699213695 3878909942 4035217385]
	//map[771555706:http://127.0.0.1:10 1060649185:http://127.0.0.1:100 1143388674:http://127.0.0.1:100 1406271808:http://127.0.0.1:1 2186562736:http://127.0.0.1:10 2217771288:http://127.0.0.1:1 2699213695:http://127.0.0.1:100 3878909942:http://127.0.0.1:10 4035217385:http://127.0.0.1:1]
	t.Log(hash)

	require.Equal(t, hash.Get("2"), "http://127.0.0.1:10")
	hash.Add("http://127.0.0.1:0")
	t.Log(hash)
	require.Equal(t, hash.Get("2"), "http://127.0.0.1:0")
}
