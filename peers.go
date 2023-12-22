package easyCache

// PeerPicker 获取节点接口
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 获取缓存接口
type PeerGetter interface {
	Get(group string, string2 string) ([]byte, error)
}
