package handler

import (
	dblayer "cloud_distributed_storage/database"
	"cloud_distributed_storage/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Response struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token,omitempty"`
}

// SignupHandler: handle user signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignupHandler called")
	if r.Method == http.MethodPost {
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println("Error decoding request:", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		fmt.Printf("Received signup request: %+v\n", req)

		username := req.Username
		password := req.Password
		email := req.Email

		if len(username) < 2 || len(password) < 4 || len(email) == 0 {
			http.Error(w, fmt.Sprintf("Invalid parameter: username=%d, password=%d, email=%d", len(username), len(password), len(email)), http.StatusBadRequest)
			return
		}
		enc_password := util.Sha1([]byte(password + pwd_salt))
		suc := dblayer.UserSignup(username, enc_password, email)
		resp := Response{}
		if suc {
			resp.Code = 0
			resp.Msg = "SUCCESS"
		} else {
			resp.Code = -1
			resp.Msg = "FAILED"
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	http.Error(w, "Request error", http.StatusMethodNotAllowed)
}

// SignInHandler: login api
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		//r.ParseForm()
		//username := r.Form.Get("username")
		//password := r.Form.Get("password")
		username := req.Username
		password := req.Password
		encPassword := util.Sha1([]byte(password + pwd_salt))

		pwdChecked := dblayer.UserSignIn(username, encPassword)
		resp := Response{}
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return

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
	}
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
	username, err := ValidateTokenAndGetUsername(token)
	if err != nil {
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
