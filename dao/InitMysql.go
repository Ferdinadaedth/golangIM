package dao

import (
	"database/sql"
	"fmt"
)

const (
	userName = "root"
	Password = "yx041110"
	ip       = "newk8s.ferdinandaedth.top"
	port     = "3306"
	dbName   = "userdb"
)

var db *sql.DB

func InitMysql() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(10) // 设置连接池中的最大连接数
	db.SetMaxIdleConns(5)  // 设置连接池中的最大空闲连接数
}
