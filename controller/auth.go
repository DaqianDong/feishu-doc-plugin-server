package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var oauthEndpoint = oauth2.Endpoint{
	AuthURL:  "https://accounts.feishu.cn/open-apis/authen/v1/authorize",
	TokenURL: "https://open.feishu.cn/open-apis/authen/v2/oauth/token",
}

var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("APP_ID"),
	ClientSecret: os.Getenv("APP_SECRET"),
	RedirectURL:  "http://localhost:8081/callback", // 请先添加该重定向 URL，配置路径：开发者后台 -> 开发配置 -> 安全设置 -> 重定向 URL -> 添加
	Endpoint:     oauthEndpoint,
	Scopes:       []string{"offline_access"}, // 如果你不需要 refresh_token，请注释掉该行，否则你需要先申请 offline_access 权限方可使用，配置路径：开发者后台 -> 开发配置 -> 权限管理
}

func IndexController(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	var username string
	session := sessions.Default(c)
	if session.Get("user") != nil {
		username = session.Get("user").(string)
	}
	html := fmt.Sprintf(`<html><head><style>body{font-family:Arial,sans-serif;background:#f4f4f4;margin:0;display:flex;justify-content:center;align-items:center;height:100vh}.container{text-align:center;background:#fff;padding:30px;border-radius:10px;box-shadow:0 0 10px rgba(0,0,0,0.1)}a{padding:10px 20px;font-size:16px;color:#fff;background:#007bff;border-radius:5px;text-decoration:none;transition:0.3s}a:hover{background:#0056b3}}</style></head><body><div class="container"><h2>欢迎%s！</h2><a href="/login">使用飞书登录</a></div></body></html>`, username)
	c.String(http.StatusOK, html)
}

func LoginController(c *gin.Context) {
	session := sessions.Default(c)

	// 生成随机 state 字符串，你也可以用其他有意义的信息来构建 state
	state := fmt.Sprintf("%d", rand.Int())
	// 将 state 存入 session 中
	session.Set("state", state)
	// 生成 PKCE 需要的 code verifier
	verifier := oauth2.GenerateVerifier()
	// 将 code verifier 存入 session 中
	session.Set("code_verifier", verifier)
	session.Save()

	url := oauthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	// 用户点击登录时，重定向到授权页面
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func OauthCallbackController(c *gin.Context) {
	session := sessions.Default(c)
	ctx := context.Background()

	// 从 session 中获取 state
	expectedState := session.Get("state")
	state := c.Query("state")

	// 如果 state 不匹配，说明是 CSRF 攻击，拒绝处理
	if state != expectedState {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", expectedState, state)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	// 如果 code 为空，说明用户拒绝了授权
	if code == "" {
		log.Printf("error: %s", c.Query("error"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	codeVerifier, _ := session.Get("code_verifier").(string)
	// 使用获取到的 code 获取 token
	token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		log.Printf("oauthConfig.Exchange() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	client := oauthConfig.Client(ctx, token)

	req, err := http.NewRequest("GET", "https://open.feishu.cn/open-apis/authen/v1/user_info", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// 使用 token 发起请求，获取用户信息
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Get() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()

	var user struct {
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("json.NewDecoder() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	// 后续可以用获取到的用户信息构建登录态，此处仅为示例，请勿直接使用
	session.Set("user", user.Data.Name)
	session.Save()

	c.Header("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<html><head><style>body{font-family:Arial,sans-serif;background:#f4f4f4;margin:0;display:flex;justify-content:center;align-items:center;height:100vh}.container{text-align:center;background:#fff;padding:30px;border-radius:10px;box-shadow:0 0 10px rgba(0,0,0,0.1)}a{padding:10px 20px;font-size:16px;color:#fff;background:#007bff;border-radius:5px;text-decoration:none;transition:0.3s}a:hover{background:#0056b3}}</style></head><body><div class="container"><h2>你好，%s！</h2><p>你已成功完成授权登录流程。</p><a href="/">返回主页</a></div></body></html>`, user.Data.Name)
	c.String(http.StatusOK, html)
}
