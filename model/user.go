package model

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type Getquestion struct {
	Questionid int    `json:"questionid"`
	Username   string `json:"username"`
	Question   string `json:"question"`
}
type Getmessage struct {
	Messageid int    `json:"messageid"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Content   string `json:"content"`
}
type Getanswer struct {
	Id         int    `json:"answerid"`
	Questionid string `json:"questionid"`
	Username   string `json:"username"`
	Answer     string `json:"answer"`
	Pid        string `json:"pid"`
}
type Question struct {
	Question string `form:"question" json:"question" binding:"required"`
}
type Answer struct {
	Questionid string `form:"questionid" json:"questionid" binding:"required"`
	Answer     string `form:"answer" json:"answer" binding:"required"`
}
type Comment struct {
	Answerid string `form:"answerid" json:"answerid" binding:"required"`
	Comment  string `form:"comment" json:"comment" binding:"required"`
}

type Message struct {
	content  string `form:"content" json:"content" binding:"required"`
	receiver string `form:"recerver" json:"receiver" binding:"required"`
}
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
