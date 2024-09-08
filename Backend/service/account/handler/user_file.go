package handler

import (
	"cloud_distributed_storage/Backend/service/account/proto"
	"context"
)

func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFiles, res *proto.ResUserFiles) error {
	return nil
}

func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.ResUserFileRename) error {
	return nil
}
