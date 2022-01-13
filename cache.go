package thqcache

import (
	"github.com/ithaiq/thqcache/lru"
	"sync"
)

//cache lru Cache封装 + 锁并发控制
type cache struct {
	lru      *lru.Cache
	mu       sync.Mutex
	maxBytes int64
}

func (this *cache) add(key string, value ByteValue) {
	this.mu.Lock()
	defer this.mu.Unlock()
	//懒加载 性能考虑
	if this.lru == nil {
		this.lru = lru.NewCache(this.maxBytes, nil)
	}
	this.lru.Add(key, value)
}

func (this *cache) get(key string) (value ByteValue, ok bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.lru == nil {
		return
	}
	if v, ok := this.lru.Get(key); ok {
		return v.(ByteValue), true
	}
	return
}
