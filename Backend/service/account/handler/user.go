package handler

import (
	"cloud_distributed_storage/Backend/common"
	cfg "cloud_distributed_storage/Backend/config"
	dblayer "cloud_distributed_storage/Backend/database"
	"cloud_distributed_storage/Backend/service/account/proto"
	"cloud_distributed_storage/Backend/util"
	"context"
	"fmt"
	"time"
)

type User struct{}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// Signup : RPC handler for user signup
func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, res *proto.ResSignup) error {
	username := req.Username
	password := req.Password
	email := req.Email

	if len(username) < 2 || len(password) < 4 || len(email) == 0 {
		res.Code = common.StatusParamInvalid
		res.Message = "Invalid parameter"
		return nil
	}
	enc_password := util.Sha1([]byte(password + cfg.Pwd_salt))
	suc := dblayer.UserSignup(username, enc_password, email)
	if suc {
		res.Code = common.StatusOK
		res.Message = "SIGNUP SUCCESS"
	} else {
		res.Code = common.StatusRegisterFailed
		res.Message = "SIGNUP FAILED"
	}
	return nil
}

// Login: RPC handler for user signin
func (u *User) Login(ctx context.Context, req *proto.ReqLogin, res *proto.ResLogin) error {

	username := req.Username
	password := req.Password
	encPassword := util.Sha1([]byte(password + cfg.Pwd_salt))

	pwdChecked := dblayer.UserSignIn(username, encPassword)

	if !pwdChecked {
		res.Code = common.StatusLoginFailed
		res.Message = "LOGIN FAILED"
	} else {
		token := GenToken(username)
		upRes := dblayer.UpdateToken(username, token)
		if !upRes {
			res.Code = common.StatusUserNotExists
			res.Message = "UPDATE FAILED"
		} else {
			res.Code = common.StatusOK
			res.Message = "LOGIN SUCCESS"
		}
	}
	return nil
}
