package model

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type Friend struct {
	Friendid   string `form:"Friendid " json:"Friendid " binding:"required"`
	Friendname string `form:"Friendname" json:"Friendname" binding:"required"`
}
type Smessage struct {
	Sender  string `form:"sender" json:"sender" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}
type Gmessage struct {
	Sender  string `form:"sender" json:"sender" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
	Groupid string `form:"groupid" json:"groupid" binding:"required"`
}
type GroupMember struct {
	MemberID   int    `json:"memberid"`
	MemberName string `json:"membername"`
	GroupID    int    `json:"groupid"`
}

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
