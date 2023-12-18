package easyCache

type ByteView struct {
	b []byte
}

// Len implement Value interface
func (v ByteView) Len() int64 {
	return int64(len(v.b))
}

func (v ByteView) String() string {
	return string(v.b)
}
