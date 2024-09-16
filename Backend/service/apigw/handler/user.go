package handler

import (
	cmn "cloud_distributed_storage/Backend/common"
	cfg "cloud_distributed_storage/Backend/config"
	userProto "cloud_distributed_storage/Backend/service/account/proto"
	dlProto "cloud_distributed_storage/Backend/service/download/proto"
	upProto "cloud_distributed_storage/Backend/service/upload/proto"
	"cloud_distributed_storage/Backend/util"
	"context"
	hystrix "github.com/asim/go-micro/plugins/wrapper/breaker/hystrix/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	ratelimit2 "github.com/juju/ratelimit"
	"log"
	"net/http"
)

var (
	userCli userProto.UserService
	upCli   upProto.UploadService
	dlCli   dlProto.DownloadService
)

func init() {
	//配置请求容量及qps
	bRate := ratelimit2.NewBucketWithRate(100, 1000)
	service := micro.NewService(
		micro.Flags(cmn.CustomFlags...),
		micro.WrapClient(ratelimit.NewClientWrapper(bRate, false)), //加入限流功能, false为不等待(超限即返回请求失败)
		micro.WrapClient(hystrix.NewClientWrapper()),               // 加入熔断功能, 处理rpc调用失败的情况(cirucuit breaker)
	)

	//init service
	service.Init()

	cli := service.Client()

	// 初始化一个account服务的客户端
	userCli = userProto.NewUserService("go.micro.service.user", cli)
	// 初始化一个upload服务的客户端
	upCli = upProto.NewUploadService("go.micro.service.upload", cli)
	// 初始化一个download服务的客户端
	dlCli = dlProto.NewDownloadService("go.micro.service.download", cli)
}

// SignupHandler: register api
func SignupHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	email := c.Request.FormValue("email")
	phone := c.Request.FormValue("phone")

	rpcResp, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: username,
		Password: password,
		Email:    email,
		Phone:    phone,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
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

	// 登录成功，返回用户信息
	cliResp := util.RespMsg{
		Code: int(cmn.StatusOK),
		Msg:  "登录成功",
		Data: struct {
			Location      string
			Username      string
			Token         string
			UploadEntry   string
			DownloadEntry string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    rpcResp.Token,
			// UploadEntry:   upEntryResp.Entry,
			// DownloadEntry: dlEntryResp.Entry,
			UploadEntry:   cfg.UploadLBHost,
			DownloadEntry: cfg.DownloadLBHost,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
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
