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
	Password = "yx041110"
	ip       = "newk8s.ferdinandaedth.top"
	port     = "3306"
	dbName   = "userdb"
)

func Creategroup(groupname string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO Groups (groupname) VALUES (?)", groupname)
	if err != nil {
		panic(err.Error())
	}
}
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
func Selectid(username string) string {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var userid string
	err = db.QueryRow("SELECT userid FROM user WHERE username=?", username).Scan(&userid)
	if err != nil {
		panic(err.Error())
	}
	return userid
}
func Addmember(id, groupid string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var membername string
	err = db.QueryRow("SELECT username FROM user WHERE userid=?", id).Scan(&membername)
	if err != nil {
		panic(err.Error())
	}
	// 插入用户记录
	_, err = db.Exec("INSERT INTO groupmember (memberid,membername,groupid) VALUES (?, ?,?)", id, membername, groupid)
	if err != nil {
		panic(err.Error())
	}
}
func Selectgroupid(groupname string) string {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var groupid string
	err = db.QueryRow("SELECT groupid FROM Groups WHERE groupname=?", groupname).Scan(&groupid)
	if err != nil {
		panic(err.Error())
	}
	return groupid
}
func Addfriend(username, friendid string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var userid int
	var uuserid int
	err = db.QueryRow("SELECT userid FROM user WHERE username=?", friendid).Scan(&userid)
	if err != nil {
		panic(err.Error())
	}
	err = db.QueryRow("SELECT userid FROM user WHERE username=?", username).Scan(&uuserid)
	if err != nil {
		panic(err.Error())
	}
	// 插入用户记录
	_, err = db.Exec("INSERT INTO user_relation (userid, friendid,friendname) VALUES (?, ?,?)", username, userid, friendid)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("INSERT INTO user_relation (friendname, friendid,userid) VALUES (?, ?,?)", username, uuserid, friendid)
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
func Getgroupmember(groupid string) []model.GroupMember {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 查询数据库
	rows, err := db.Query("SELECT memberid, membername, groupid FROM groupmember WHERE groupid = ?", groupid)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var groupMembers []model.GroupMember

	// 遍历查询结果，将数据存储到切片中
	for rows.Next() {
		var member model.GroupMember
		err := rows.Scan(&member.MemberID, &member.MemberName, &member.GroupID)
		if err != nil {
			panic(err.Error())
		}
		groupMembers = append(groupMembers, member)
	}

	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
	return groupMembers
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
