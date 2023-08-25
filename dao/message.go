package dao

import (
	"database/sql"
	"fmt"
)

func DepositSmessages(sender, receiver, message, mtype string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO smessage (sender,receiver,message,type) VALUES (?,?,?,?)", sender, receiver, message, mtype)
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
	_, err = db.Exec("INSERT INTO gmessage (sender,groupid,message,type) VALUES (?,?,?,?)", sender, groupid, message, mtype)
	if err != nil {
		panic(err.Error())
	}
}
