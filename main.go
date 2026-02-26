package main

import (
	"fmt"
	"log"
	"oauth-test/controller"
	"oauth-test/infra/larkclient"
	"oauth-test/infra/ocr"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	larkclient.Init(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"))
	ocr.Init(os.Getenv("OCR_KEY"))

	// 端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()
	r.Use(Cors())

	// 使用 Cookie 存储 session
	store := cookie.NewStore([]byte("secret")) // 此处仅为示例，务必不要硬编码密钥
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", controller.IndexController)
	r.GET("/login", controller.LoginController)
	r.GET("/callback", controller.OauthCallbackController)
	r.GET("/whiteboard", controller.Wrap(controller.WhiteboardController))

	fmt.Println("Server running on http://localhost:" + port)
	log.Fatal(r.Run(":" + port))
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 1. 设置允许的 Header 字段
		c.Header("Access-Control-Allow-Origin", "*") // * 表示允许所有来源，生产环境建议指定具体域名
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 2. 处理浏览器发送的 OPTIONS 预检请求
		if method == "OPTIONS" {
			c.AbortWithStatus(204) // 返回 204 No Content
			return
		}

		c.Next()
	}
}
