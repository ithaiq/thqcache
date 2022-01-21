# thqcache
一个简单可用的分布式缓存系统

* commit学习步骤

1. [LRU缓存淘汰算法封装实现](https://github.com/ithaiq/thqcache/commit/c77df27674c02e0f9fe37a5e17be57e31babcc5a)
2. [Cache并发控制封装LRU+缓存组Group实现+数据格式封装](https://github.com/ithaiq/thqcache/commit/912f6dcc33cbf186ed0d2473fff3e188327eeb91)
3. [Http服务端封装](https://github.com/ithaiq/thqcache/commit/5e305badba77e15502f04c44f385e0ec0405b445)
4. [一致性哈希算法封装](https://github.com/ithaiq/thqcache/commit/e95039a283a6078487e811f4e351196eb9df7480)
5. [HTTP客户端封装&完善整个缓存流程](https://github.com/ithaiq/thqcache/commit/941e5e6fb9b38b6b201a6a2bc6230d1b1c1a235d)
6. [解决并发导致缓存击穿](https://github.com/ithaiq/thqcache/commit/9898e98495d846301a14d4e29ea94b4a57ba9e69)
7. [使用protobuf优化节点通信](https://github.com/ithaiq/thqcache/commit/ed6334d98f5dfd0391ad4a0720017ccc90df32d9)

* 调用关系链如下
  ![调用关系](https://github.com/ithaiq/thqcache/raw/master/gocall.png)