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
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
