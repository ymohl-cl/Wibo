//header
package usr

import (
	"container/list"
	//"fmt"
	"github.com/Wibo/src/db"
)

/*
** Log is a last signal to device
 */

type User struct {
	Device  int64
	Login   string
	id_user int64
	mail    string
	Log     int
}

type All_users struct {
	Lst_users list.List
}

func (User *User) User_is_online() bool {
	if User.Log == 0 {
		return true
	} else {
		return false
	}
}

func (Lst_users *All_users) Del_user(del_user *User) {
	return
}

func (Lst_users *All_users) Add_new_user(new_user *User) {
	return
}

func (Lst_users *All_users) Print_users() {
	// Print All_users
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (Lst_users *All_users) initUser(uid int64, login string, mail string) *User {
	t := new(User)
	t.id_user = uid
	t.Login = login
	t.mail = mail
	return (t)
}

// Get all users
func (Lst_users *All_users) Get_users() error {

	var err error
	//	lUser := All_users{}
	rows, err := db.Db.Query("SELECT id_user, login, mail FROM \"user\";")
	for rows.Next() {
		var idUser int64
		var login string
		var mail string
		err = rows.Scan(&idUser, &login, &mail)
		checkErr(err)
		//lUser.List.PushBack(lUser.initUser(idUser, login, mail))
	}
	//	Lst_users lUser
	return nil
}
