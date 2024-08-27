package config

import (
	cmn "cloud_distributed_storage/Backend/common"
)

const (
	// TempLocalRootDir : 本地临时存储地址的路径
	TempLocalRootDir = "/data/fileserver/"
	// TempPartRootDir : 分块文件在本地临时存储地址的路径
	TempPartRootDir = "/data/fileserver_part/"
	// CephRootDir : Ceph的存储路径prefix
	CephRootDir = "/ceph"
	// OSSRootDir : OSS的存储路径prefix
	S3RootDir = "S3/"
	// CurrentStoreType : 设置当前文件的存储类型
	CurrentStoreType = cmn.StoreLocal
)
