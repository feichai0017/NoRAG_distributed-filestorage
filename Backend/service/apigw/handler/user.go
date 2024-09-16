package handler

import (
	cmn "cloud_distributed_storage/Backend/common"
	userProto "cloud_distributed_storage/Backend/service/account/proto"
	"cloud_distributed_storage/Backend/util"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

var (
	userCli userProto.UserService
)

func init() {
	service := micro.NewService()

	//init service
	service.Init()

	//init user service rpc client
	userCli = userProto.NewUserService("go.micro.service.user", service.Client())
}

// SignupHandler: register api
func SignupHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	email := c.Request.FormValue("email")

	rpcResp, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: username,
		Password: password,
		Email:    email,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": rpcResp.Code,
		"msg":  rpcResp.Message,
	})
}

// SignInHandler: login api
func SignInHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	rpcResp, err := userCli.Login(context.TODO(), &userProto.ReqLogin{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if rpcResp.Code != cmn.StatusOK {
		c.JSON(200, gin.H{
			"msg":  "login failed",
			"code": rpcResp.Code,
		})
		return
	}
}

// SignOutHandler: 处理登出请求
func SignOutHandler(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username not found in context"})
		return
	}
	_, err := userCli.Logout(context.TODO(), &userProto.ReqLogout{
		Token: token.(string),
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
}

// DeleteUserHandler: 处理删除用户请求
func DeleteUserHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	_, err := userCli.DeleteAccount(context.TODO(), &userProto.ReqDeleteAccount{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username not found in context"})
		return
	}

	resp, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: username.(string),
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3. 组装并且响应用户数据
	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: gin.H{
			"Username": username,
			"SignupAt": resp.SignupAt,
			// TODO: 完善其他字段信息
			"LastActive": resp.LastActiveAt,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}
