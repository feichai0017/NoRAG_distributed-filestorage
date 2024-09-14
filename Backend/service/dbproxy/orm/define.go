package orm

import (
	"database/sql"
	"time"
)

// TableFile 文件表结构
type TableFile struct {
	ID       int64
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	OwnerID  int64
	CreateAt time.Time
	UpdateAt time.Time
	Status   int
}

// TableUser 用户表结构
type TableUser struct {
	ID             int64
	Username       string
	UserPwd        string
	Email          string
	Phone          string
	EmailValidated bool
	PhoneValidated bool
	SignupAt       time.Time
	LastActive     time.Time
	Profile        sql.NullString
	Status         int
}

// TableRole 角色表结构
type TableRole struct {
	ID          int64
	RoleName    string
	Description sql.NullString
	CreateAt    time.Time
	UpdateAt    time.Time
}

// TableUserRole 用户角色关联表结构
type TableUserRole struct {
	ID       int64
	UserID   int64
	RoleID   int64
	CreateAt time.Time
}

// TablePermission 权限表结构
type TablePermission struct {
	ID         int64
	RoleID     sql.NullInt64
	UserID     sql.NullInt64
	FileID     int64
	PermRead   bool
	PermWrite  bool
	PermDelete bool
	PermShare  bool
	ExpireTime sql.NullTime
	CreateAt   time.Time
	UpdateAt   time.Time
}

// TableUserFile 用户文件表结构
type TableUserFile struct {
	ID          int64
	UserID      int64
	FileID      int64
	FileName    string
	FileHash    string
	FileSize    int64
	UploadAt    time.Time
	LastUpdated time.Time
	Status      int
}

// ExecResult 执行结果
type ExecResult struct {
	Suc  bool        `json:"suc"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
