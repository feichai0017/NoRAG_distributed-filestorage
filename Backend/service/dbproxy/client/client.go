package client

import (
	"cloud_distributed_storage/Backend/service/dbproxy/orm"
	dbProto "cloud_distributed_storage/Backend/service/dbproxy/proto"
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt time.Time
	OwnerID  int64
}

var (
	dbCli dbProto.DBProxyService
)

func Init(service micro.Service) {
	// 初始化 dbproxy service
	dbCli = dbProto.NewDBProxyService("go.micro.service.dbproxy", service.Client())
}

func TableFileToFileMeta(tfile orm.TableFile) FileMeta {
	return FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
		UploadAt: tfile.CreateAt,
		OwnerID:  tfile.OwnerID,
	}
}

// execAction : send request to dbproxy to execute action
func execAction(funcName string, paramJson []byte) (*dbProto.ResExec, error) {
	return dbCli.ExecuteAction(context.TODO(), &dbProto.ReqExec{
		Actions: []*dbProto.SingleAction{
			&dbProto.SingleAction{
				Name:   funcName,
				Params: paramJson,
			},
		},
	})
}

// parseBody : parse response rpc body
func parseBody(res *dbProto.ResExec) *orm.ExecResult {
	if res == nil || res.Data == nil {
		log.Println("parseBody: res or res.Data is nil")
		return nil
	}
	resList := []orm.ExecResult{}
	err := json.Unmarshal(res.Data, &resList)
	if err != nil {
		log.Printf("parseBody: json.Unmarshal failed, err:%v", err)
		return nil
	}

	if len(resList) > 0 {
		return &resList[0]
	}
	return nil
}

func ToTableUser(src interface{}) orm.TableUser {
	user := orm.TableUser{}
	_ = mapstructure.Decode(src, &user)
	return user
}

func ToTableFile(src interface{}) orm.TableFile {
	file := orm.TableFile{}
	_ = mapstructure.Decode(src, &file)
	return file
}

func ToTableFiles(src interface{}) []orm.TableFile {
	var files []orm.TableFile
	_ = mapstructure.Decode(src, &files)
	return files
}

func ToTableUserFile(src interface{}) orm.TableUserFile {
	userFile := orm.TableUserFile{}
	_ = mapstructure.Decode(src, &userFile)
	return userFile
}

func ToTableUserFiles(src interface{}) []orm.TableUserFile {
	var userFiles []orm.TableUserFile
	_ = mapstructure.Decode(src, &userFiles)
	return userFiles
}

func GetFileMeta(filehash string) (*orm.ExecResult, error) {
	uInfo, err := json.Marshal([]string{filehash})
	res, err := execAction("/file/GetFileMeta", uInfo)
	if err != nil {
		return nil, err
	}
	return parseBody(res), nil
}

func GetFileMetaList(limit int) (*orm.ExecResult, error) {
	uInfo, err := json.Marshal([]int{limit})
	res, err := execAction("/file/GetFileMetaList", uInfo)
	if err != nil {
		return nil, err
	}
	return parseBody(res), nil
}

// OnFileUploadFinished : when file upload finished, save file meta to db
func OnFileUploadFinished(fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location, fmeta.OwnerID})
	res, err := execAction("/file/OnFileUploadFinished", uInfo)
	return parseBody(res), err
}

func UpdateFileLocation(filehash, location string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash, location})
	res, err := execAction("/file/UpdateFileLocation", uInfo)
	return parseBody(res), err
}

func UserSignup(username, encPasswd, email, phone string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd, email, phone})
	res, err := execAction("/user/UserSignup", uInfo)
	return parseBody(res), err
}

func UserLogin(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := execAction("/user/UserLogin", uInfo)
	return parseBody(res), err
}

func GetUserInfo(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/GetUserInfo", uInfo)
	return parseBody(res), err
}

func UserExist(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/UserExist", uInfo)
	return parseBody(res), err
}

func UpdateToken(username, token string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, token})
	res, err := execAction("/user/UpdateToken", uInfo)
	return parseBody(res), err
}

func QueryUserFileMeta(userID int64, filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID, filehash})
	res, err := execAction("/ufile/QueryUserFileMeta", uInfo)
	return parseBody(res), err
}

func QueryUserFileMetas(userID int64, limit int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID, limit})
	res, err := execAction("/ufile/QueryUserFileMetas", uInfo)
	return parseBody(res), err
}

// OnUserFileUploadFinished : 新增/更新文件元信息到mysql中
func OnUserFileUploadFinished(userID int64, fileID int64, fileName string, fileSize int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID, fileID, fileName, fileSize})
	res, err := execAction("/ufile/OnUserFileUploadFinished", uInfo)
	return parseBody(res), err
}

func RenameFileName(userID int64, filehash, filename string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID, filehash, filename})
	res, err := execAction("/ufile/RenameFileName", uInfo)
	return parseBody(res), err
}

// 新增的函数

func CreateRole(roleName string, description string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{roleName, description})
	res, err := execAction("/role/CreateRole", uInfo)
	return parseBody(res), err
}

func GetRoleInfo(roleID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{roleID})
	res, err := execAction("/role/GetRoleInfo", uInfo)
	return parseBody(res), err
}

func UpdateRole(roleID int64, roleName, description string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{roleID, roleName, description})
	res, err := execAction("/role/UpdateRole", uInfo)
	return parseBody(res), err
}

func DeleteRole(roleID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{roleID})
	res, err := execAction("/role/DeleteRole", uInfo)
	return parseBody(res), err
}

func ListRoles() (*orm.ExecResult, error) {
	res, err := execAction("/role/ListRoles", nil)
	return parseBody(res), err
}

func AssignRoleToUser(userID, roleID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID, roleID})
	res, err := execAction("/role/AssignRoleToUser", uInfo)
	return parseBody(res), err
}

func RemoveRoleFromUser(userID, roleID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID, roleID})
	res, err := execAction("/role/RemoveRoleFromUser", uInfo)
	return parseBody(res), err
}

func GetUserRoles(userID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID})
	res, err := execAction("/role/GetUserRoles", uInfo)
	return parseBody(res), err
}

func GrantPermission(roleID, userID, fileID int64, permRead, permWrite, permDelete, permShare bool, expireTime *time.Time) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{roleID, userID, fileID, permRead, permWrite, permDelete, permShare, expireTime})
	res, err := execAction("/permission/GrantPermission", uInfo)
	return parseBody(res), err
}

func RevokePermission(roleID, userID, fileID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{roleID, userID, fileID})
	res, err := execAction("/permission/RevokePermission", uInfo)
	return parseBody(res), err
}

func CheckPermission(userID, fileID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID, fileID})
	res, err := execAction("/permission/CheckPermission", uInfo)
	return parseBody(res), err
}

func ListUserPermissions(userID int64) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{userID})
	res, err := execAction("/permission/ListUserPermissions", uInfo)
	return parseBody(res), err
}
