package orm

import (
	"database/sql"
)

// TableFile 文件表结构
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	Status   int
	Ext1     int
	Ext2     sql.NullString
}

// TableUser 用户表结构
type TableUser struct {
	UserName   string
	UserPwd    string
	Email      string
	Phone      string
	SignupAt   string
	LastActive string
	Profile    sql.NullString
	Status     int
}

// TableRole 角色表结构
type TableRole struct {
	RoleName    string
	Description sql.NullString
	CreateAt    string
	UpdateAt    string
}

// TableUserRole 用户角色关联表结构
type TableUserRole struct {
	UserName string
	RoleName string
	CreateAt string
}

// TablePermission 权限表结构
type TablePermission struct {
	RoleName   sql.NullString
	UserName   sql.NullString
	FileSha1   string
	PermRead   bool
	PermWrite  bool
	PermDelete bool
	PermShare  bool
	ExpireTime sql.NullTime
	CreateAt   string
	UpdateAt   string
}

// TableUserFile 用户文件表结构
type TableUserFile struct {
	UserName       string
	FileHash       string
	FileName       string
	FileSize       int64
	UploadAt       string
	LastUpdated    string
	Status         int
	DowndloadCount int
}

// ExecResult 执行结果
type ExecResult struct {
	Suc  bool        `json:"suc"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
