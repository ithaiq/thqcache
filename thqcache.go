package thqcache

import (
	"fmt"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

//接口型函数
//方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数。

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//Group 负责与外部交互控制主流程
type Group struct {
	name   string
	getter Getter
	c      cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		c:      cache{maxBytes: maxBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (this *Group) Add(key string, value ByteValue) {
	this.c.add(key, value)
}

func (this *Group) Get(key string) (ByteValue, error) {
	if key == "" {
		return ByteValue{}, fmt.Errorf("key is empty")
	}
	if v, ok := this.c.get(key); ok {
		return v, nil
	}
	return this.loadLocal(key)
}

//loadLocal 加载本地数据到缓存
func (this *Group) loadLocal(key string) (ByteValue, error) {
	bytes, err := this.getter.Get(key)
	if err != nil {
		return ByteValue{}, err
	}
	value := ByteValue{b: cloneBytes(bytes)}
	this.Add(key, value)
	return value, nil
}
