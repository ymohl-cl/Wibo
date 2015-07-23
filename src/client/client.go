package client

import (
	"container/list"
	"database/sql"
	"fmt"
	//"github.com/Wibo/src/db"
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
	Db, err := sql.Open("postgres", "user=wibo  password='wibo' dbname=wibo_base sslmode=disable host=localhost port=49155")
	rows, err := Db.Query("SELECT id_user, login, mail FROM \"user\";")
	if err != nil {
		return
	}
	for rows.Next() {
		var id_user int64
		var login string
		var mail string
		err = rows.Scan(&id_user, &login, &mail)

		fmt.Println("idUser | username |  mail ")
		fmt.Printf("+%v %v %v", id_user, login, mail)
	}
	return
}
