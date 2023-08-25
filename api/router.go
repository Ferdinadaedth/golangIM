package api

import (
	"github.com/gin-gonic/gin"
	"golangIM/api/middleware"
	"golangIM/dao"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/verify", middleware.JWTAuthMiddleware())
	r.POST("/register", register) // 注册
	r.POST("/login", login)       // 登录
	r.POST("/addfriend", middleware.JWTAuthMiddleware(), addfriend)
	r.POST("/deletefriend", middleware.JWTAuthMiddleware(), deletefriend)
	r.GET("/ws/:userID", middleware.JWTAuthMiddleware(), websocketHandler)
	r.GET("/groupws/:groupID", middleware.JWTAuthMiddleware(), groupWebSocketHandler)
	r.GET("/upload/:userID", middleware.JWTAuthMiddleware(), uploadImage)
	r.GET("/groupload/:groupID", middleware.JWTAuthMiddleware(), groupuploadImage)
	r.GET("/creategroup/:groupname", middleware.JWTAuthMiddleware(), creategroup)
	r.POST("/invitefriend", invitefriend)
	r.GET("/getfriend", middleware.JWTAuthMiddleware(), dao.Getfriend)
	r.GET("/getsmessage", middleware.JWTAuthMiddleware(), getsmessage)
	r.GET("/getgmessage", middleware.JWTAuthMiddleware(), getgmessage)
	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.POST("/get", getUsernameFromToken)
	}

	r.Run(":8088") // 跑在 8088 端口上
}
