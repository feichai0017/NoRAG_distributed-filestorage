package handler

import (
	dblayer "cloud_distributed_storage/database"
	"cloud_distributed_storage/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler: handle user signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}
	enc_password := util.Sha1([]byte(password + pwd_salt))
	suc := dblayer.UserSignup(username, enc_password)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// SignInHandler: login api
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/login.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPassword := util.Sha1([]byte(password + pwd_salt))

	pwdChecked := dblayer.UserSignIn(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	//w.Write([]byte("http://"+r.Host+))
	//resp := util.RespMsg{
	//	Code: 0,
	//	Msg:  "OK",
	//	Data: struct {
	//		Location string
	//		Username string
	//		Token    string
	//	}{
	//		Location: "http://" + r.Host + "/service",
	//		Username: username,
	//		Token:    token,
	//	},
	//}
	http.Redirect(w, r, "/service", http.StatusFound)
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("/usr/local/Distributed_system/cloud_distributed_storage/static/view/userinfo.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	isTokenValid := IsTokenValid(token)
	if !isTokenValid {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

func IsTokenValid(token string) bool {
	if len(token) <= 8 {
		return false // Token 长度不足
	}
	tokenTimestamp := token[len(token)-8:]

	// 将时间戳转换为 int64 类型
	ts, err := fmt.Sscanf(tokenTimestamp, "%x", new(int64))
	if err != nil {
		fmt.Println("Invalid token timestamp:", err)
		return false
	}

	// 检查 token 是否在 2 小时内有效
	tokenTime := time.Unix(int64(ts), 0)
	if time.Since(tokenTime) > 2*time.Hour {
		return false // Token 已过期
	}
	return true
}
