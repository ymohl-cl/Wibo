package users

import (
	"Wibo/protocol"
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"strings"
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

type Coordinate struct {
	Lon float64
	Lat float64
}

type StatsUser struct {
	CreationDate  time.Time /* Date de creation set par le serveur */
	NbrBallCreate int64     /* Nombre de ballon cree par l'user */
	NbrCatch      int64     /* Nombre de ballon catche par l'user */
	NbrSend       int64     /* Nombre de ballon envoye par l'user */
	NbrFollow     int64     /* Nombre de ballon Follow par l'user */
	NbrMessage    int64     /* Nombre de message ecris par l'user */
}

type User struct {
	Id          int64      /* Id bdd de l'user set par la base de donnee */
	Mail        string     /* Mail de l'user */
	NbrBallSend int        /** Nombre de ballon envoye le meme jour */
	Coord       Coordinate /* Interface Coord pour connaitre la position de l'user */
	Log         time.Time  /** Date of the last query doned by user: Peut devenir deprecated */
	Followed    *list.List /* Value: *list.Element.Value.(*ballon.Ball) */
	Possessed   *list.List /* Value: *list.Element.Value.(*ballon.Ball) */
	HistoricReq *list.List /* Liste d'interface History, compose l'historique des requetes utilisateurs */
	Stats       *StatsUser /* Interface Stats, Statistique de la vie du ballon */
	/* Les valeurs suivantes sont deprecated */
	//Device      *list.List /* Value: Device */
	//Password    string // pas utile car la comparaison sera faite avec la bdd
	//Login string Not use, not login
}

type All_users struct {
	Ulist *list.List
	//Id_max int64
	GlobalStat *StatsUser /* Stats globale a tous les utilisateur de WIbo */
	NbrUsers   int64
}

type userError struct {
	prob string
	err  error
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

func FoundUserOnListLvl2(lst *list.List, email string) *list.Element {
	euser := lst.Front()
	user := euser.Value.(*list.Element).Value.(*User)

	for euser != nil && strings.Compare(user.Mail, email) != 0 {
		euser = euser.Next()
		user = euser.Value.(*list.Element).Value.(*User)
	}
	if euser != nil {
		return euser.Value.(*list.Element)
	}
	return nil
}

func FoundUserOnListLvl1(lst *list.List, email string) *list.Element {
	euser := lst.Front()
	user := euser.Value.(*User)

	for euser != nil && strings.Compare(user.Mail, email) != 0 {
		euser = euser.Next()
		user = euser.Value.(*User)
	}
	if euser != nil {
		return euser
	}
	return nil
}

/*
** Search request' user on list parameter, if not found, search in all list.
** If not found, return nil, else, check request' password.
** If Password is OK return user else return nil
 */
func (ulist *All_users) Check_user(request *list.Element, Db *sql.DB, History *list.List) *list.Element {
	req := request.Value.(protocol.Request)
	user := FoundUserOnListLvl2(History, req.Spec.(protocol.Log).Email)
	if user == nil {
		user = FoundUserOnListLvl1(ulist.Ulist, req.Spec.(protocol.Log).Email)
	}
	if user != nil {
		user = CheckPasswordUser(user, req.Spec.(protocol.Log).Pswd, Db)
	}
	return user
}

/******************************************************************************/
/********************************* MERGE JAIME ********************************/
/******************************************************************************/

func CheckValidMail(email string) bool {
	tmp := valid.IsEmail(email)
	if tmp == true {
		fmt.Println("Email ok:")
	} else {
		fmt.Println("Email KO:")
	}
	fmt.Println(email)
	return tmp
}

func CheckPasswordUser(user *list.Element, pass string, Db *sql.DB) *list.Element {
	var err error
	passb := []byte(pass)
	rows, err := Db.Query("SELECT id_user, mail, bpass FROM \"user\" WHERE id_user=$1;", user.Value.(*User).Id)
	for rows.Next() {
		var idUser int64
		var mailq string
		var bpass []byte
		err = rows.Scan(&idUser, &mailq, bpass)
		checkErr(err)
		if bcrypt.CompareHashAndPassword(bpass, passb) != nil {
			return nil
		}
	}
	return user
}

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

func (e *userError) Error() string {
	return fmt.Sprintf("%s - %v", e.prob, e.err)
}

/**
* Insert new user to wibo_base
*	constrain valid_mail(text) default verification
* query to check if mail is already registered
**/

func (Lst_users *All_users) Add_new_user(new_user *User, Db *sql.DB, Pass string) (bool, error) {
	var err error

	if len(Pass) > 1 {
		/* really danger below */
		Pass = "ThisIsAPasswordDefault2015OP"
	}
	bpass, err := bcrypt.GenerateFromPassword([]byte(Pass), 15)
	if err != nil {
		return false, &userError{"Error add new user", err}
	}

	if len(new_user.Mail) > 0 {
		if valid.IsEmail(new_user.Mail) != true {
			return false, errors.New("Wrong mail format")
		}
	}
	_, err = Db.Query("INSERT INTO \"user\" (id_type_g, groupname, passbyte, lastlogin, creationdate, mail) VALUES ($1, $2, $3, $4, $5, $6);", 1, "particulier", bpass, time.Now(), time.Now(), new_user.Mail)

	//	_, err = Db.Query("INSERT INTO \"user\" (id_type_g, groupname, passbyte, creationdate, mail) VALUES ($1, $2, $3, $4, $5);", 1, "particulier", new_user.Mail, bpass, time.Now(), time.Now(), new_user.Mail)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
Insert user default
	insert a user with default data
*/
func (Lst_users *All_users) AddNewDefaultUser(Db *sql.DB) *list.Element {
	bpass, err := bcrypt.GenerateFromPassword([]byte("Password_default2015"), 15)
	if err != nil {
		fmt.Println("Error crypt")
		fmt.Println(err)
		return nil
	}
	rows, err := Db.Query(
		"INSERT INTO \"user\" (id_type_g, groupname, passbyte, lastlogin, creationdate, mail) VALUES ($1, $2, $3, $4, $5, make_uid()) RETURNING id_user;", 2, "user_default", bpass, time.Now(), time.Now())
	if err != nil {
		fmt.Println("Error insert")
		fmt.Println(err)
		return nil
	}
	for rows.Next() {
		var IdUserDefault int64
		err = rows.Scan(&IdUserDefault)
		Lst_users.Ulist.PushBack(&User{Id: IdUserDefault, Followed: list.New(), Possessed: list.New(), HistoricReq: list.New()})
		//Lst_users.Ulist.PushBack(&User{Login: "user_default", Id: IdUserDefault})
	}
	return Lst_users.Ulist.Back()
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
		//		var login string
		var mailq string
		//		err = rows.Scan(&idUser, &login, &mailq)
		err = rows.Scan(&idUser, &mailq)
		checkErr(err)
		//		return initUser(idUser, login, mailq)
		return initUser(idUser, mailq)
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
		//		fmt.Printf("%v | %v | %v \n", e.Value.(*User).Id, e.Value.(*User).Login, e.Value.(*User).Mail)
		fmt.Printf("%v | %v \n", e.Value.(*User).Id, e.Value.(*User).Mail)
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
//func initUser(uid int64, login string, mail string) *User {
func initUser(uid int64, mail string) *User {
	t := new(User)
	//t.Login = login
	t.Id = uid
	t.Mail = mail
	return (t)
}

/**
* NewUser
* Create User instance and return it
**/
func (Lusr *All_users) NewUser(mail string) *User {
	return &User{Mail: mail}
}

//func (Lusr *All_users) NewUser(login string, mail string, pass string) *User {
//	return &User{Login: login, Mail: mail, Password: pass}
//}

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
	//	rows, err := Db.Query("SELECT id_user, login, mail, password FROM \"user\";")
	rows, err := Db.Query("SELECT id_user, mail, password FROM \"user\";")
	checkErr(err)
	for rows.Next() {
		var idUser int64
		//		var login string
		var mailq string
		var pass string
		//		err = rows.Scan(&idUser, &login, &mailq, &pass)
		err = rows.Scan(&idUser, &mailq, &pass)
		checkErr(err)
		//		lDevice := Lusr.GetDevicesByIdUser(idUser, Db)
		//		lUser.PushBack(&User{Login: login, Id: idUser, Mail: mailq, Device: lDevice, Followed: list.New()})
		//		lUser.PushBack(&User{Id: idUser, Mail: mailq, Device: lDevice, Followed: list.New()})
		lUser.PushBack(&User{Id: idUser, Mail: mailq, Followed: list.New()})
	}
	Lusr.Ulist.Init()
	Lusr.Ulist.PushFrontList(lUser)
	return nil
}
