package thqcache

import "github.com/ithaiq/thqcache/proto"

//PeerPicker 节点选择器
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//PeerGetter 从节点获取数据
type PeerGetter interface {
	Get(in *proto.Request, out *proto.Response) error
}
