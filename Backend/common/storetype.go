package common

type StoreType int

const (
	Local StoreType = iota
	Ceph
	S3
)

func (s StoreType) String() string {
	return [...]string{"Local", "Ceph", "S3"}[s]
}
