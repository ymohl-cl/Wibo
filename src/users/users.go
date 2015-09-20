package users

import (
	"container/list"
	"database/sql"
	"fmt"
	"protocol"
	"strconv"
	//	_ "github.com/lib/pq"
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
	IdMobile    int64      /* type int64 is temporary */
	History_req *list.List /* Value: History */
}

type User struct {
	Id       int64
	Login    string
	Mail     string
	Password string
	Device   *list.List /* Value: Device */
	Log      time.Time  /*Date of the last query */
	Followed *list.List /* Value: *list.Element.Value.(*ballon.Ball) */
}

type All_users struct {
	Ulist  *list.List
	Id_max int64
}

func (User *User) User_is_online() bool {
	t_now := time.Now()
	t_user := User.Log
	if t_user.Hour() == t_now.Hour() && t_user.Minute() > t_now.Minute()-2 {
		return true
	} else {
		return false
	}
}

/*
** Manage users's connexion
 */
func (ulist *All_users) Check_user(request *list.Element, Db *sql.DB) (user *list.Element, err error) {
	user = ulist.Ulist.Front()
	var device *list.Element

	rqt := request.Value.(protocol.Request)
	for user != nil {
		device = user.Value.(*User).Device.Front()
		for device != nil && device.Value.(Device).IdMobile != rqt.Deviceid {
			device = device.Next()
		}
		if device != nil {
			break
		}
		user = user.Next()
	}
	if user == nil {
		usr := new(User)
		var hist_device Device
		usr.Device = list.New()
		usr.Log = time.Now()
		usr.Followed = list.New()
		hist_device.IdMobile = request.Value.(protocol.Request).Deviceid
		hist_device.History_req = list.New()
		hist_device.History_req.PushFront(History{time.Now(), request.Value.(protocol.Request).Rtype})
		usr.Device.PushFront(hist_device)
		user = ulist.Ulist.PushBack(usr)
		ulist.Add_new_user(usr, Db)
	} else {
		device.Value.(Device).History_req.PushFront(History{time.Now(), request.Value.(protocol.Request).Rtype})
		user.Value.(*User).Log = time.Now()
	}
	return user, nil
}

/******************************************************************************/
/********************************* MERGE JAIME ********************************/
/******************************************************************************/

/**
* Delete user from id and mail
*	TODO: del device
**/
func (Lst_users *All_users) Del_user(del_user *User, Db *sql.DB) (executed bool, err error) {
	stm, err := Db.Prepare("DELETE FROM  \"user\" WHERE id_user=$1")
	_, err = stm.Exec(del_user.Id)
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
	for e := LstU.Ulist.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v | %v | %v \n", e.Value.(*User).Id, e.Value.(*User).Login, e.Value.(*User).Mail)
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
	t.Id = uid
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
		v, err := strconv.Atoi(idDevice)
		checkErr(err)
		lDevice.PushBack(&Device{IdMobile: int64(v), History_req: list.New()})
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
		lUser.PushBack(&User{Login: login, Id: idUser, Mail: mailq, Device: lDevice, Followed: list.New()})
	}
	Lusr.Ulist.Init()
	Lusr.Ulist.PushFrontList(lUser)
	return nil
}
