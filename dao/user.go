package dao

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golangIM/model"
	"golangIM/utils"
	"log"
	"net/http"
)

const (
	userName = "root"
	Password = "123456"
	ip       = "newk8s.ferdinandaedth.top"
	port     = "3306"
	dbName   = "userdb"
)

func SelectUser(username string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 查询用户名是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE username=?", username).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		return true
	} else {
		return false
	}

}
func Addfriend(username, friendid string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var userid int

	err = db.QueryRow("SELECT userid FROM user WHERE username=?", friendid).Scan(&userid)
	if err != nil {
		panic(err.Error())
	}
	// 插入用户记录
	_, err = db.Exec("INSERT INTO user_relation (userid, friendid,friendname) VALUES (?, ?,?)", username, userid, friendid)
	if err != nil {
		panic(err.Error())
	}
}
func Deletefriend(username, friendid string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 插入用户记录
	_, err = db.Exec("DELETE FROM user_relation WHERE userid = 'username' AND friendid = 'friendid'")
	if err != nil {
		panic(err.Error())
	}
}

// AddUser 添加用户
func AddUser(username, password string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 插入用户记录
	_, err = db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		panic(err.Error())
	}
}
func SelectPasswordFromUsername(username string) string {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 查询密码
	var password string
	err = db.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&password)
	if err != nil {
		panic(err.Error())
	}
	return password
}
func Getfriend(c *gin.Context) {
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "username not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var friends []model.Friend
	rows, errq := db.Query("select friendid,friendname from user_relation where userid = ?", username)
	if errq != nil {
		log.Fatal(errq.Error)
		return
	}
	//遍历结果
	for rows.Next() {
		var u model.Friend
		errn := rows.Scan(&u.Friendid, &u.Friendname)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		friends = append(friends, u)
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})

}
