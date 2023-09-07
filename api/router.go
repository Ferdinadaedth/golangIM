package api

import (
	"github.com/gin-gonic/gin"
	"golangIM/api/middleware"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/verify", middleware.JWTAuthMiddleware())
	r.POST("/register", register) // 注册
	r.POST("/login", login)       // 登录
	r.POST("/addfriend", middleware.JWTAuthMiddleware(), addfriend)
	r.GET("/getre", middleware.JWTAuthMiddleware(), getfriendre)
	r.GET("/handlere/:reid/:result", middleware.JWTAuthMiddleware(), processfriendre)
	r.GET("/getgroups", middleware.JWTAuthMiddleware(), getgroup)
	r.POST("/deletefriend", middleware.JWTAuthMiddleware(), deletefriend)
	r.GET("/ws/:userID", middleware.JWTAuthMiddleware(), sendsmessage)
	r.GET("/groupws/:groupID", middleware.JWTAuthMiddleware(), sendgmessage)
	r.GET("/upload/:userID", middleware.JWTAuthMiddleware(), uploadImage)
	r.GET("/groupload/:groupID", middleware.JWTAuthMiddleware(), groupuploadImage)
	r.GET("/creategroup/:groupname", middleware.JWTAuthMiddleware(), creategroup)
	r.POST("/invitefriend", middleware.JWTAuthMiddleware(), invitefriend)
	r.GET("/getfriends", middleware.JWTAuthMiddleware(), getfriend)
	r.GET("/getsmessage", middleware.JWTAuthMiddleware(), getsmessage)
	r.GET("/getgmessage", middleware.JWTAuthMiddleware(), getgmessage)
	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.POST("/get", getUsernameFromToken)
	}

	r.Run(":8088") // 跑在 8088 端口上
}
