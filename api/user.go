package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golangIM/api/middleware"
	"golangIM/dao"
	"golangIM/model"
	"golangIM/utils"
	"net/http"
	"time"
)

func getgroup(c *gin.Context) {
	username := dao.Getusername(c)
	groups := dao.Getgroup(username)
	c.JSON(http.StatusOK, gin.H{
		"status": "200",
		"groups": groups,
	})
}
func getfriend(c *gin.Context) {
	username := dao.Getusername(c)
	friends := dao.Getfriend(username)
	c.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"friends": friends,
	})
}
func getfriendre(c *gin.Context) {
	username := dao.Getusername(c)
	friendre := dao.GetFriendre(username)
	c.JSON(http.StatusOK, gin.H{
		"status":   "200",
		"friendre": friendre,
	})
}
func processfriendre(c *gin.Context) {
	reid := c.Param("reid")
	result := c.Param("result")
	dao.Handlere(reid, result)
	utils.RespSuccess(c, "处理成功")
}
func register(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 验证用户名是否重复
	flag := dao.SelectUser(username)
	fmt.Println(flag)
	if flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user already exists")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Println(err)
	}
	encodePWD := string(hash)
	dao.AddUser(username, encodePWD)
	// 以 JSON 格式返回信息
	utils.RespSuccess(c, "add user successful")
}
func login(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否存在
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}

	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username)
	// 若不正确则传出错误
	err := bcrypt.CompareHashAndPassword([]byte(selectPassword), []byte(password)) //验证（对比）
	if err != nil {
		utils.RespFail(c, "wrong password")
		return
	}
	c.SetCookie("gin_demo_cookie", "test", 3600, "/", "localhost", false, true)
	// 创建一个我们自己的声明
	claim := model.MyClaims{
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 过期时间
			Issuer:    "ferdinand",                          // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, _ := token.SignedString(middleware.Secret)
	c.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "登录成功",
		"token":   tokenString,
	})
}
func addfriend(c *gin.Context) {
	friendid := c.PostForm("friendid")
	username := dao.Getusername(c)
	dao.Addfriend(username, friendid)
	utils.RespSuccess(c, "成功发送好友请求")
}
func deletefriend(c *gin.Context) {
	friendid := c.PostForm("friendid")
	username := dao.Getusername(c)
	dao.Deletefriend(username, friendid)
	utils.RespSuccess(c, "成功删除好友")
}
func getUsernameFromToken(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"status":   "200",
		"username": username,
	})
}
