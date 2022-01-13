package lru

import "container/list"

//Cache LRU map+双向链表
type Cache struct {
	cache     map[string]*list.Element
	maxBytes  int64
	nowBytes  int64
	ll        *list.List
	onRemoved func(key string, value Value)
}

//Element 节点
type Element struct {
	key   string
	value Value
}

//Value 返回节点占用内存大小
type Value interface {
	Len() int
}

func NewCache(maxBytes int64, onRemoved func(key string, value Value)) *Cache {
	return &Cache{
		cache:     make(map[string]*list.Element),
		maxBytes:  maxBytes,
		ll:        list.New(),
		onRemoved: onRemoved,
	}
}

func (this *Cache) Add(key string, value Value) {
	if element, ok := this.cache[key]; ok {
		this.ll.MoveToFront(element)
		kv := element.Value.(*Element)
		kv.value = value
		this.nowBytes += int64(value.Len()) - int64(kv.value.Len())
	} else {
		element = this.ll.PushFront(&Element{key, value})
		this.cache[key] = element
		this.nowBytes += int64(len(key)) + int64(value.Len())
	}
	//淘汰最近最少使用的元素
	for this.maxBytes != 0 && this.maxBytes < this.nowBytes {
		this.RemoveOldest()
	}
}

func (this *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := this.cache[key]; ok {
		this.ll.MoveToFront(element)
		kv := element.Value.(*Element)
		return kv.value, true
	}
	return
}

func (this *Cache) RemoveOldest() {
	element := this.ll.Back()
	if element != nil {
		this.ll.Remove(element)
		kv := element.Value.(*Element)
		delete(this.cache, kv.key)
		this.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if this.onRemoved != nil {
			this.onRemoved(kv.key, kv.value)
		}
	}
}

func (this *Cache) Len() int {
	return this.ll.Len()
}
