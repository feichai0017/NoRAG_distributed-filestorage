package handler

import (
	dblayer "cloud_distributed_storage/Backend/database"
	"cloud_distributed_storage/Backend/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token,omitempty"`
}

// SignupHandler: handle user signup
func SignupHandler(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	username := req.Username
	password := req.Password
	email := req.Email

	if len(username) < 2 || len(password) < 4 || len(email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter"})
		return
	}
	enc_password := util.Sha1([]byte(password + pwd_salt))
	suc := dblayer.UserSignup(username, enc_password, email)
	resp := UserResponse{}
	if suc {
		resp.Code = 0
		resp.Msg = "SUCCESS"
	} else {
		resp.Code = -1
		resp.Msg = "FAILED"
	}
	c.JSON(http.StatusOK, resp)
}

// SignInHandler: login api
func SignInHandler(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	username := req.Username
	password := req.Password
	encPassword := util.Sha1([]byte(password + pwd_salt))

	pwdChecked := dblayer.UserSignIn(username, encPassword)
	resp := UserResponse{}
	if !pwdChecked {
		resp.Code = -1
		resp.Msg = "FAILED"
	} else {
		token := GenToken(username)
		upRes := dblayer.UpdateToken(username, token)
		if !upRes {
			resp.Code = -1
			resp.Msg = "FAILED"
		} else {
			resp.Code = 0
			resp.Msg = "SUCCESS"
			resp.Token = token
		}
	}
	c.JSON(http.StatusOK, resp)
}

func UserInfoHandler(c *gin.Context) {
	username, err := GetUsernameFromContext(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "User not found"})
	}
	token := c.Query("token")
	username, err = ValidateTokenAndGetUsername(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid token"})
		return
	}

	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.JSON(http.StatusOK, resp)
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

func ValidateTokenAndGetUsername(token string) (string, error) {
	if len(token) <= 8 {
		return "", fmt.Errorf("invalid token length")
	}

	// 验证 token 时效性
	tokenTimestamp := token[len(token)-8:]
	var ts int64
	_, err := fmt.Sscanf(tokenTimestamp, "%x", &ts)
	if err != nil {
		return "", fmt.Errorf("invalid token timestamp")
	}

	tokenTime := time.Unix(ts, 0)
	if time.Since(tokenTime) > 2*time.Hour {
		return "", fmt.Errorf("token expired")
	}
	username := dblayer.QueryUserByToken(token)
	if len(username) == 0 {
		return "", fmt.Errorf("username not found: %v", err)
	}
	return username, nil
}
