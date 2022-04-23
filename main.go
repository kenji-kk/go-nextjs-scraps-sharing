package main

import (
	"fmt"
	"hobby/db"
	"hobby/handler"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db.DbConnect()
}

func main() {
        fmt.Println("起動")
        db.JwtSetup()
        
	r := gin.Default()
		r.Use(cors.New(cors.Config{
			// 許可したいHTTPメソッドの一覧
			AllowMethods: []string{
					"POST",
					"GET",
					"OPTIONS",
					"PUT",
					"DELETE",
			},
			// 許可したいHTTPリクエストヘッダの一覧
			AllowHeaders: []string{
					"Access-Control-Allow-Headers",
					"Content-Type",
					"Content-Length",
					"Accept-Encoding",
					"X-CSRF-Token",
					"Authorization",
			},
			// 許可したいアクセス元の一覧
			AllowOrigins: []string{
					"*",
			},
			// 自分で許可するしないの処理を書きたい場合は、以下のように書くこともできる
			// AllowOriginFunc: func(origin string) bool {
			//  return origin == "https://www.example.com:8080"
			// },
			// preflight requestで許可した後の接続可能時間
			// https://godoc.org/github.com/gin-contrib/cors#Config の中のコメントに詳細あり
			MaxAge: 24 * time.Hour,
	}))
        
        r.POST("/signup", handler.Signup)
        r.POST("/signin", handler.Signin)

        r.Run(":8080")
}
