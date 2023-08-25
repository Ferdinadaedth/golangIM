package dao

import (
	"database/sql"
	"fmt"
	"golangIM/model"
)

func ProcessGMessages(username string) []model.Gmessage {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// 查询 groupmember 表中指定 username 的行数
	var count int
	var gmessages []model.Gmessage
	err = db.QueryRow("SELECT COUNT(*) FROM groupmember WHERE membername=?", username).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	if count > 0 {
		// 查询 groupmember 表中满足 username 的所有记录
		rows, err := db.Query("SELECT groupid FROM groupmember WHERE membername=?", username)
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()

		// 创建一个切片用于存储所有的消息记录

		// 遍历查询结果
		for rows.Next() {
			var groupID string
			if err := rows.Scan(&groupID); err != nil {
				panic(err.Error())
			}

			// 查询 gmessage 表中指定 groupid 的消息记录
			groupMessages, err := getGroupMessagesFromDB(groupID)
			if err != nil {
				panic(err.Error())
			}

			// 将 groupMessages 添加到 gmessages 切片中
			gmessages = append(gmessages, groupMessages...)
		}
	}
	return gmessages
}

func getGroupMessagesFromDB(groupID string) ([]model.Gmessage, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// 查询 gmessage 表中指定 groupid 的消息记录
	rows, err := db.Query("SELECT sender, content,groupid FROM gmessage WHERE groupid=?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 创建一个切片用于存储消息记录
	var messages []model.Gmessage

	// 遍历查询结果
	for rows.Next() {
		var sender, content, groupid string
		if err := rows.Scan(&sender, &content, &groupid); err != nil {
			return nil, err
		}

		// 创建 Gmessage 结构体并添加到切片中
		message := model.Gmessage{Sender: sender, Content: content, Groupid: groupid}
		messages = append(messages, message)
	}

	return messages, nil
}

func Getsmessage(username string) []model.Smessage {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 查询数据库
	rows, err := db.Query("SELECT content,sender FROM smessage WHERE  receiver= ?", username)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var smessages []model.Smessage

	// 遍历查询结果，将数据存储到切片中
	for rows.Next() {
		var smess model.Smessage
		err := rows.Scan(&smess.Content, &smess.Sender)
		if err != nil {
			panic(err.Error())
		}
		smessages = append(smessages, smess)
	}

	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
	return smessages
}
func DepositSmessages(sender, receiver, message, mtype string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO smessage (sender,receiver,content,type) VALUES (?,?,?,?)", sender, receiver, message, mtype)
	if err != nil {
		panic(err.Error())
	}
}
func DepositGmessages(sender, groupid, message, mtype string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO gmessage (sender,groupid,content,type) VALUES (?,?,?,?)", sender, groupid, message, mtype)
	if err != nil {
		panic(err.Error())
	}
}
