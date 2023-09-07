package dao

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golangIM/cache"
	"golangIM/model"
	"log"
)

func Getgroup(username string) []model.Group {
	var groups []model.Group
	rows, errq := db.Query("select groupid from groupmember where membername = ? ", username)
	if errq != nil {
		log.Fatal(errq.Error)
	}
	//遍历结果
	for rows.Next() {
		var u model.Group
		errn := rows.Scan(&u.Groupid)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		groups = append(groups, u)
	}
	return groups
}
func GetFriendre(username string) []model.Friendre {
	var friendre []model.Friendre
	defaultstatus := "0"
	rows, errq := db.Query("select friendid,friendname,reid from user_relation where username = ? AND status =?", username, defaultstatus)
	if errq != nil {
		log.Fatal(errq.Error)
	}
	//遍历结果
	for rows.Next() {
		var u model.Friendre
		errn := rows.Scan(&u.Friendid, &u.Friendname, &u.Reid)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		friendre = append(friendre, u)
	}
	return friendre
}
func Handlere(reid, result string) {
	_, err := db.Exec("UPDATE user_relation SET status = ? WHERE reid = ? ", result, reid)
	if err != nil {
		panic(err.Error())
	}
	var username string
	var friendname string
	err = db.QueryRow("SELECT username,friendname FROM user_relation WHERE reid=?", reid).Scan(&username, &friendname)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("UPDATE user_relation SET status = ? WHERE username = ? AND friendname = ? ", result, friendname, username)
	if err != nil {
		panic(err.Error())
	}
}
func Creategroup(groupname string) {
	_, err := db.Exec("INSERT INTO Groups (groupname) VALUES (?)", groupname)
	if err != nil {
		panic(err.Error())
	}
}
func SelectUser(username string) bool {
	// 查询用户名是否存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE username=?", username).Scan(&count)
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
	var userid string
	err := db.QueryRow("SELECT userid FROM user WHERE username=?", username).Scan(&userid)
	if err != nil {
		panic(err.Error())
	}
	return userid
}
func Selectusername(id string) string {
	var username string
	err := db.QueryRow("SELECT username FROM user WHERE userid=?", id).Scan(&username)
	if err != nil {
		panic(err.Error())
	}
	return username
}
func Addmember(id, groupid string) {
	var membername string
	err := db.QueryRow("SELECT username FROM user WHERE userid=?", id).Scan(&membername)
	if err != nil {
		panic(err.Error())
	}
	// 插入用户记录
	_, err = db.Exec("INSERT INTO groupmember (memberid,membername,groupid) VALUES (?, ?,?)", id, membername, groupid)
	if err != nil {
		panic(err.Error())
	}
	err = cache.DeleteCache("groupmembers" + groupid)
	if err != nil {
		log.Fatal(err.Error)
	}
}
func Selectgroupid(groupname string) string {
	var groupid string
	err := db.QueryRow("SELECT groupid FROM Groups WHERE groupname=?", groupname).Scan(&groupid)
	if err != nil {
		panic(err.Error())
	}
	return groupid
}
func Addfriend(username, friendid string) {
	friendname := Selectusername(friendid)
	userid := Selectid(username)
	defaultstatus := "0"
	// 插入用户记录
	_, err := db.Exec("INSERT INTO user_relation (username, friendid,friendname,status) VALUES (?, ?,?,?)", username, friendid, friendname, defaultstatus)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("INSERT INTO user_relation (username, friendid,friendname,status) VALUES (?, ?,?,?)", friendname, userid, username, defaultstatus)
	if err != nil {
		panic(err.Error())
	}
}
func Deletefriend(username, friendid string) {

	// 插入用户记录
	_, err := db.Exec("DELETE FROM user_relation WHERE userid = ? AND friendid = ?", username, friendid)
	if err != nil {
		panic(err.Error())
	}
}

// AddUser 添加用户
func AddUser(username, password string) {

	// 插入用户记录
	_, err := db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		panic(err.Error())
	}
}
func SelectPasswordFromUsername(username string) string {

	// 查询密码
	var password string
	err := db.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&password)
	if err != nil {
		panic(err.Error())
	}
	return password
}
func Getgroupmember(groupid string) []model.GroupMember {
	var groupMembers []model.GroupMember
	// 查询数据库
	err := cache.GetCache("groupmembers", &groupMembers)
	if err == nil {
		return groupMembers
	}
	rows, err := db.Query("SELECT memberid, membername, groupid FROM groupmember WHERE groupid = ?", groupid)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

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
	err = cache.SetCache("groupmembers"+groupid, &groupMembers)
	if err != nil {
		log.Fatal(err.Error)
	}
	return groupMembers
}
func Getfriend(username string) []model.Friend {
	var friends []model.Friend
	newstatus := "1"
	rows, errq := db.Query("select friendid,friendname from user_relation where username = ? AND status =?", username, newstatus)
	if errq != nil {
		log.Fatal(errq.Error)
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
	return friends

}
func Getusername(c *gin.Context) string {
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		fmt.Printf("username not found")
	}
	username, ok := value.(string)
	if !ok {
		fmt.Printf("username is not a string")
	}
	return username
}
