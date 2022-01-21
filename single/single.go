package single

import "sync"

//缓存雪崩：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。
//
//缓存击穿：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。
//
//缓存穿透：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。

//目的所有用户都能收到结果，请求是在服务端阻塞的，等待某一个查询返回结果的，其余请求直接复用这个结果了
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

//为何不用channel通道？
//如果用信道，接受和发送需要一一对应 waitgroup 有 Add(1) 和 Done() 是一一对应的，但是可以有多个请求同时调用 Wait()，同时等待该任务结束， 一般锁和信道是做不到这一点的。
func (this *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	this.mu.Lock()
	if this.m == nil {
		this.m = make(map[string]*call)
	}
	if c, ok := this.m[key]; ok {
		this.mu.Unlock()
		c.wg.Wait() // 如果请求正在进行中，则等待
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	this.m[key] = c
	this.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()

	return c.val, c.err
}
