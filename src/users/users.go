package users

import (
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

/**
** Date est la date a laquelle la requete a ete effectue.
** Type_req_client et le type de requete effectue.
**/
type History struct {
	Date            time.Time
	Type_req_client int16
}

/**
** -type Device
** IdMobile est l'identifiant unique du mobile.
** Pour le moment le format exact de l'IdMobile est inconnu.
** History_req est une liste qui sera l'historique des requetes du client
** depuis ce device.
**/

type Device struct {
	IdMobile    int64
	History_req *list.List
}

type User struct {
	Device    *list.List
	Login     string
	Id_user   int64
	Mail      string
	Password  string
	Log       time.Time
	LastLogin string
}

type All_users struct {
	Lst_users *list.List
}

/* Definis si l'utilisateur est considere en ligne ou pas avec un timeout de 2 min */
func (User *User) User_is_online() bool {
	t_now := time.Now()
	t_user := User.Log
	if t_user.Hour() == t_now.Hour() && t_user.Minute() > t_now.Minute()-2 {
		return true
	} else {
		return false
	}
}

/**
* Delete user from id and mail
*	TODO: del device
**/
func (Lst_users *All_users) Del_user(del_user *User, Db *sql.DB) (executed bool, err error) {
	stm, err := Db.Prepare("DELETE FROM  \"user\" WHERE id_user=$1")
	_, err = stm.Exec(del_user.Id_user)
	checkErr(err)
	executed = true
	return executed, err
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

/**
* SelectUser
* Create an instance of an User with their data with some Id
* return an instance User
 */

func (LstU *All_users) SelectUser(idUser int64, Db *sql.DB) *User {
	var err error
	rows, err := Db.Query("SELECT id_user, login, mail FROM \"user\" WHERE id_user=$1;", idUser)
	for rows.Next() {
		var idUser int64
		var login string
		var mailq string
		err = rows.Scan(&idUser, &login, &mailq)
		checkErr(err)
		return initUser(idUser, login, mailq)
	}
	return nil
}

/**
* Print_users
* Print Value.Interface.Element to output
 */

func (LstU *All_users) Print_users() {
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
* GetDevicesByIdUser
* Query function getDevicesUserIdi  with prototype:
*	FUNCTION getDevicesByUserId(iduser integer) RETURNS TABLE(macaddr varchar(18))
* Return a pointer on new Device list created
**/
func (Lusr *All_users) GetDevicesByIdUser(idUser int64, Db *sql.DB) *list.List {

	lDevice := list.New()
	stm, err := Db.Prepare("SELECT getDevicesByUserId($1)")
	checkErr(err)
	rows, err := stm.Query(idUser)
	for rows.Next() {
		var idDevice string
		err = rows.Scan(&idDevice)
		checkErr(err)
		lDevice.PushBack(idDevice)
	}
	return lDevice
}

/**
* Get_users
* Query the user table join device and create new *listList Pointer
**/

func (Lusr *All_users) Get_users(Db *sql.DB) error {

	var err error
	lUser := list.New()
	rows, err := Db.Query("SELECT id_user, login, mail, password FROM \"user\";")
	checkErr(err)
	for rows.Next() {
		var idUser int64
		var login string
		var mailq string
		var pass string
		err = rows.Scan(&idUser, &login, &mailq, &pass)
		checkErr(err)
		lDevice := Lusr.GetDevicesByIdUser(idUser, Db)
		lUser.PushBack(User{Login: login, Id_user: idUser, Mail: mailq, Device: lDevice})
	}
	Lusr.Lst_users = Lusr.Lst_users.Init()
	Lusr.Lst_users.PushFrontList(lUser)
	return nil
}
