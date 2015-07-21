package client

import (
	"container/list"
	"fmt"
	"github.com/Wibo/src/db"
)

type User struct {
	Device int64
	Log    int
}

type All_users struct {
	Lst_users list.List
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Get_users() {
	// Get all users
	var err error
	Db, err := db.OpenCo(err)
	rows, err := Db.Query("SELECT id_user, login, mail FROM \"user\";")
	if err != nil {
		return
	}
	for rows.Next() {
		var id_user int64
		var login string
		var mail string
		err = rows.Scan(&id_user, &login, &mail)
		checkErr(err)

		fmt.Println("idUser | username |  mail ")
		fmt.Println(id_user, login, mail)
	}
	return rows
}
