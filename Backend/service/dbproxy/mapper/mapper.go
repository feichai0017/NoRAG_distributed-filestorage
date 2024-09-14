package mapper

import (
	"cloud_distributed_storage/Backend/service/dbproxy/orm"
	"fmt"
	"reflect"
)

var funcs = map[string]interface{}{
	"/file/OnFileUploadFinished": orm.OnFileUploadFinished,
	"/file/GetFileMeta":          orm.GetFileMeta,
	"/file/GetFileMetaList":      orm.GetFileMetaList,
	"/file/UpdateFileLocation":   orm.UpdateFileLocation,

	"/user/UserSignup":        orm.UserSignup,
	"/user/UserLogin":         orm.UserLogin,
	"/user/UserExist":         orm.UserExist,
	"/user/UpdateToken":       orm.UpdateToken,
	"/user/GetUserInfo":       orm.GetUserInfo,
	"/user/UserLogout":        orm.UserLogout,
	"/user/DeleteUserAccount": orm.DeleteUserAccount,

	"/ufile/OnUserFileUploadFinished": orm.OnUserFileUploadFinished,
	"/ufile/QueryUserFileMetas":       orm.QueryUserFileMetas,
	"/ufile/QueryUserFileMeta":        orm.QueryUserFileMeta,
	"/ufile/UpdateUserFileName":       orm.RenameFileName,
	"/ufile/DeleteUserFile":           orm.DeleteUserFile,

	// 新增的RBAC相关函数映射
	"/role/CreateRole":             orm.CreateRole,
	"/role/GetRoleInfo":            orm.GetRoleInfo,
	"/user/AssignRole":             orm.UpdateRole,
	"/user/RemoveRole":             orm.DeleteRole,
	"/permission/GrantPermission":  orm.GrantPermission,
	"/permission/RevokePermission": orm.RevokePermission,
	"/permission/CheckPermission":  orm.CheckPermission,
}

func FunCall(name string, params ...interface{}) (result []reflect.Value, err error) {
	f, ok := funcs[name]
	if !ok {
		err = fmt.Errorf("func %s not found", name)
		return
	}
	// use reflect to call the function
	fv := reflect.ValueOf(f)
	if len(params) != fv.Type().NumIn() {
		err = fmt.Errorf("func %s need %d params", name, fv.Type().NumIn())
		return
	}
	// construct a slice of reflect.Value
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = fv.Call(in)
	return
}
