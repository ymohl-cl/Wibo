//header
package usr

import (
	"container/list"
	//"database/sql"
	"fmt"
	"github.com/Wibo/src/db"
	_ "github.com/lib/pq"
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
	Lst_users  *list.List
	next, prev *All_users
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

func (LstU *All_users) Print_users() {
	// Print All_users
	i := 0
	for e := LstU.Lst_users.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v \n", e.Value.(User).mail)
		i++
	}
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initUser(uid int64, login string, mail string) *User {
	t := new(User)
	t.id_user = uid
	t.Login = login
	t.mail = mail
	return (t)
}

// Get all users
func (Lusr *All_users) Get_users() error {

	var err error
	lUser := list.New()
	//Db, err := sql.Open("postgres", "user=wibo  password='wibo' dbname=wibo_base sslmode=disable host=localhost port=49155")
	checkErr(err)
	fmt.Printf("%v ", db.Env)
	/*rows, err := db.Db.Query("SELECT id_user, login, mail FROM \"user\";")
	checkErr(err)
	for rows.Next() {
		var idUser int64
		var login string
		var mailq string
		err = rows.Scan(&idUser, &login, &mailq)
		checkErr(err)
		lUser.PushBack(User{Login: login, id_user: idUser, mail: mailq})
	}
	Lusr.Lst_users.PushFrontList(lUser)
	*/return nil
}
