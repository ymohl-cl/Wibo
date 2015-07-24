//header
package usr

import (
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

/*
** Log is a last signal to device
 */

type User struct {
	Device    int64
	Login     string
	Id_user   int64
	Mail      string
	Log       int
	Password  string
	LastLogin string
}

type All_users struct {
	Lst_users *list.List
}

func (User *User) User_is_online() bool {
	if User.Log == 0 {
		return true
	} else {
		return false
	}
}

/**
* Delete user from id and mail
*	TODO: del device
**/
func (Lst_users *All_users) Del_user(del_user *User) {
	var err error
	tblname := "user"
	_, err = Db.Exec(
		fmt.Sprint("DELETE FROM  \"%s\" WHERE user.id_user=$1 HAVING user.mail=$2", tblname),
		del_user.Id_user, del_user.Mail)
	checkErr(err)
	return
}

/**
* Insert new user to wibo_base
*	TODO: imput verification
**/

func (Lst_users *All_users) Add_new_user(new_user *User, Db *sql.DB) {
	var err error
	tblname := "user"
	_, err = Db.Exec(
		fmt.Sprintf(
			"INSERT INTO \"%s\"(id_type_g, groupname, login, password, salt, lastlogin, creationdate, mail) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", tblname),
		1, "particulier", new_user.Login, new_user.Password, "saltTest", time.Now(), time.Now(), new_user.Mail)
	checkErr(err)
	return
}

func (LstU *All_users) Select_user(idUser int64, Db *sql.DB) *User {
	var err error
	tblname := "user"
	row, err = Db.Exec(fmt.Sprintf("SELECT user.id_user, user.login, user.lastlogin, user.creationdate, user.mail FROM \"%s\" WHERE user.id_user=$1", tblname), idUser)
	for rows.Next() {
		var idUser int64
		var login string
		var mailq string
		var creationdate string
		var lastlogin string
		err = rows.Scan(&idUser, &login, &mailq, &creationdate, &lastlogin)
		checkErr(err)
		return initUser(login, idUser, mailq)
	}
}

/**
* Print_users
* Print Value.Interface.Element to output
 */

func (LstU *All_users) Print_users() {
	// Print All_users
	i := 0
	for e := LstU.Lst_users.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v | %v | %v \n", e.Value.(User).Id_user, e.Value.(User).Login, e.Value.(User).Mail)
		i++
	}
	return
}

/**
* CheckErr
* Verify err value to stop execution by panic
**/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/**
* InitUser
* initialisation of new user instance with custom data
 */
func initUser(uid int64, login string, mail string) *User {
	t := new(User)
	t.Id_user = uid
	t.Login = login
	t.Mail = mail
	return (t)
}

/**
* NewUser
* Create User instance and return it
**/
func (Lusr *All_users) NewUser(login string, mail string, pass string) *User {
	return &User{Login: login, Mail: mail, Password: pass}
}

/**
* Get_users
* Query the user table join device and create new *listList Pointer
* TODO join table Device
 */

func (Lusr *All_users) Get_users(Db *sql.DB) error {

	var err error
	lUser := list.New()
	rows, err := Db.Query("SELECT id_user, login, mail FROM \"user\";")
	checkErr(err)
	for rows.Next() {
		var idUser int64
		var login string
		var mailq string
		err = rows.Scan(&idUser, &login, &mailq)
		checkErr(err)
		lUser.PushBack(User{Login: login, Id_user: idUser, Mail: mailq})
	}
	Lusr.Lst_users.PushFrontList(lUser)
	return nil
}
