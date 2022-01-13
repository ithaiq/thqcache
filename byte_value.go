package thqcache

//ByteValue 数据格式封装
type ByteValue struct {
	b []byte
}

func (v ByteValue) Len() int {
	return len(v.b)
}

func (v ByteValue) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteValue) String() string {
	return string(v.b)
}

//cloneBytes 深拷贝防止修改b
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
