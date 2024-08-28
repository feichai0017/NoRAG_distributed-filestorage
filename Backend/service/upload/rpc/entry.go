package rpc

import (
	"cloud_distributed_storage/Backend/service/upload/config"
	uploadProto "cloud_distributed_storage/Backend/service/upload/proto"
	"context"
)

type Upload struct {
}

func (u *Upload) UploadEntry(ctx context.Context, req *uploadProto.ReqEntry, res *uploadProto.ResEntry) error {
	res.Entry = config.UploadEntry
	return nil
}
