package handler

import (
	"cloud_distributed_storage/Backend/common"
	"cloud_distributed_storage/Backend/service/account/proto"
	dbcli "cloud_distributed_storage/Backend/service/dbproxy/client"
	"context"
	"encoding/json"
)

func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFiles, res *proto.ResUserFiles) error {
	dbResp, err := dbcli.QueryUserFileMetas(req.Username, int(req.Limit))
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return err
	}
	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}

func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.ResUserFileRename) error {
	dbResp, err := dbcli.RenameFileName(req.Username, req.Filehash, req.NewFileName)
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return err
	}

	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}
