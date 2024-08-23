package mq

import "cloud_distributed_storage/common"

// TransferData : 定义一个数据传输的结构体
type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType common.StoreType
}
