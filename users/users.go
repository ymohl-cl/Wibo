package users

import (
	"Wibo/db"
	"Wibo/protocol"
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
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
	Magnet      time.Time
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
	LogUser    *userError
}

type userError struct {
	Prob string
	Err  error
	Logf *log.Logger
}

func (User *User) MagnetisValid() bool {
	if time.Since(User.Magnet) > (60 * time.Minute) {
		return true
	}
	return false
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
	req := request.Value.(*protocol.Request)
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
func (lu *All_users) Get_GlobalStat(base *db.Env) (er error) {
	rows, err := base.Db.Query("SELECT num_users, num_follow, num_message, num_send, num_cont FROM globalStats;")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()
	rows.Scan(&lu.NbrUsers, &lu.GlobalStat.NbrFollow, &lu.GlobalStat.NbrMessage, &lu.GlobalStat.NbrSend, &lu.GlobalStat.NbrBallCreate)
	return err
}

// FUNCTION updatelocationuser(iduser integer, latitudec double precision, longitudec double precision)
// FUNCTION public.updateuser(iduser integer, latitudec double precision, longitudec double precision, log date)

func (lu *All_users) Update_users(base *db.Env) (err error) {
	u := lu.Ulist.Front()
	for u != nil {
		cu := u.Value.(*User)
		fmt.Printf("Update user %v \n", cu.Id)
		fmt.Printf("Update user %v \n", cu.Coord.Lon)
		fmt.Printf("Update user %v \n", cu.Coord.Lat)
		trow, err := base.Db.Query("SELECT updateuser($1, $2, $3, $4);", cu.Id, cu.Coord.Lon, cu.Coord.Lat, cu.Log)

		ex := lu.SetStatsByUser(cu.Id, cu.Stats, base.Db)
		if ex != true {
			fmt.Println("Fail to update user stats")
		}
		if err != nil {
			fmt.Println(err)
		}
		trow.Close()
		u = u.Next()
	}
	return err
}

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
	rows, err := Db.Query("SELECT id_user, mail, passbyte FROM \"user\" WHERE id_user=$1;", user.Value.(*User).Id)
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
	defer rows.Close()

	for rows.Next() {
		var idUser int64
		var mailq string
		var bpass []byte
		err = rows.Scan(&idUser, &mailq, &bpass)
		if err != nil {
			fmt.Println(err)
			fmt.Println(err)
			return nil
		}
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
	if err != nil {
		Lst_users.LogUser.Prob = "Del user fail"
		Lst_users.LogUser.Err = err
		return false, Lst_users.LogUser
	}
	executed = true
	return executed, err
}

func (e *userError) Error() string {
	if e.Err != nil {
		e.Logf.Println(e.Err)
	}
	return fmt.Sprintf("%s - %v", e.Prob, e.Err)
}

/**
* Insert new user to wibo_base
*	constrain valid_mail(text) default verification
* query to check if mail is already registered
	setsuserdata2(idtypeg integer,
	groupnamec character varying,
	latc double precision, lonc double precision,
	creation date,
	lastlog date,
	mailc character varying,
	pass bytea)
**/

func (Lst_users *All_users) Add_new_user(new_user *User, Db *sql.DB, Pass string) (bool, error) {
	var err error

	if len(Pass) == 0 {
		/* really danger below */
		Pass = "ThisIsAPasswordDefault2015OP"
	}
	bpass, err := bcrypt.GenerateFromPassword([]byte(Pass), 15)
	if err != nil {
		Lst_users.LogUser.Prob = "Add new user fail"
		Lst_users.LogUser.Err = err
		return false, Lst_users.LogUser
	}

	if len(new_user.Mail) > 0 {
		if valid.IsEmail(new_user.Mail) != true {
			return false, errors.New("Wrong mail format")
		}
	}
	/* set id*/
	//rows, err := Db.Query("INSERT INTO \"user\" (id_type_g, groupname, passbyte, lastlogin, creationdate, mail) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_user;", 1, "particulier", bpass, time.Now(), new_user.Stats.CreationDate, new_user.Mail)
	err = Db.QueryRow("SELECT  setsuserdata2($1, $2, $3, $4, $5, $6, $7, $8);",
		1,
		"user_particulier",
		new_user.Coord.Lat,
		new_user.Coord.Lon,
		new_user.Stats.CreationDate,
		new_user.Log,
		new_user.Mail,
		bpass).Scan(&new_user.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
Insert user default
	insert a user with default data
*/
func (Lst_users *All_users) AddNewDefaultUser(Db *sql.DB, req *protocol.Request) *list.Element {
	bpass, err := bcrypt.GenerateFromPassword([]byte("Password_default2015"), 15)
	if err != nil {
		Lst_users.LogUser.Prob = "GetDevicesByIdUser query"
		Lst_users.LogUser.Err = err
		fmt.Println(Lst_users.LogUser.Error())
		return nil
	}
	tmpUser := new(User)
	tmpUser.Log = time.Now()
	tmpUser.Followed = list.New()
	tmpUser.Possessed = list.New()
	tmpUser.HistoricReq = list.New()
	tmpUser.Coord.Lon = req.Coord.Lon
	tmpUser.Coord.Lat = req.Coord.Lat
	tmpUser.Stats = new(StatsUser)
	tmpUser.Stats.CreationDate = time.Now()
	tmpUser.Stats.NbrBallCreate = 0
	tmpUser.Stats.NbrCatch = 0
	tmpUser.Stats.NbrSend = 0
	tmpUser.Stats.NbrFollow = 0
	tmpUser.Stats.NbrMessage = 0
	//rows, err := Db.Query(
	//	"INSERT INTO \"user\" (id_type_g, groupname, passbyte, lastlogin, creationdate, mail) VALUES ($1, $2, $3, $4, $5, make_uid()) RETURNING id_user;", 2, "user_default", bpass, tmpUser.Log, tmpUser.Stats.CreationDate)
	err = Db.QueryRow("SELECT setsdefaultuserdata($1, $2, $3, $4, $5);", tmpUser.Coord.Lat, tmpUser.Coord.Lon, tmpUser.Log, tmpUser.Stats.CreationDate, bpass).Scan(&tmpUser.Id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	Lst_users.Ulist.PushBack(tmpUser)
	return Lst_users.Ulist.Back()
}

/**
* SelectUser
* Create an instance of an User with their data with some Id
* return an instance User
 */

func (LstU *All_users) SelectUser(idUser int64, Db *sql.DB) *User {
	var err error
	rows, err := Db.Query("SELECT id_user, mail FROM \"user\" WHERE id_user=$1;", idUser)
	if err != nil {
		fmt.Println(err)
		LstU.LogUser.Prob = "GetDevicesByIdUser query"
		LstU.LogUser.Err = err
		fmt.Printf(LstU.LogUser.Error())
		os.Exit(-1)
	}
	defer rows.Close()

	for rows.Next() {
		var idUser int64
		var mailq string
		err = rows.Scan(&idUser, &mailq)
		if err != nil {
			LstU.LogUser.Prob = "Select User fail"
			LstU.LogUser.Err = err
		}
		return initUser(idUser, mailq)
	}
	return nil
}

/**
* Print_users
* Print Value.Interface.Element to output
 */

func (LstU *All_users) Print_users() {
	for e := LstU.Ulist.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v | %v \n", e.Value.(*User).Id, e.Value.(*User).Mail)
	}
	return
}

/**
* InitUser
* initialisation of new user instance with custom data
 */
//func initUser(uid int64, login string, mail string) *User {
func initUser(uid int64, mail string) *User {
	t := new(User)
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

/**
* GetDevicesByIdUser
* Query function getDevicesUserIdi  with prototype:
*	FUNCTION getDevicesByUserId(iduser integer) RETURNS TABLE(macaddr varchar(18))
* Return a pointer on new Device list created
**/
func (Lusr *All_users) GetDevicesByIdUser(idUser int64, Db *sql.DB) *list.List {

	lDevice := list.New()
	stm, err := Db.Prepare("SELECT getDevicesByUserId($1)")
	if err != nil {
		Lusr.LogUser.Prob = "GetDevicesByIdUser query"
		Lusr.LogUser.Err = err
		fmt.Printf(Lusr.LogUser.Error())
	}
	rows, err := stm.Query(idUser)
	for rows.Next() {
		var idDevice string
		err = rows.Scan(&idDevice)
		if err != nil {
			Lusr.LogUser.Prob = "GetDevicesByIdUser Query rows fail"
			Lusr.LogUser.Err = err
			fmt.Printf(Lusr.LogUser.Error())
		}
		lDevice.PushBack(idDevice)
	}
	return lDevice
}

func GetCoord(position string) Coordinate {
	// Return true if 'value' char.
	f := func(c rune) bool {
		return c == '(' || c == '(' || c == ')' || c == '"' ||
			c == 'P' || c == 'O' || c == 'I' || c == 'N' ||
			c == 'T'
	}
	// Separate into fields with func.
	fields := strings.FieldsFunc(position, f)
	// Separate into cordinates  with Fields.
	point := strings.Fields(fields[0])
	long, _ := strconv.ParseFloat(point[0], 6)
	lat, _ := strconv.ParseFloat(point[1], 6)
	return Coordinate{Lon: long, Lat: lat}
}

/**
* Get_users
* Query the user table join device and create new *listList Pointer
**/
func (Lusr *All_users) Get_users(Db *sql.DB, fl *log.Logger) error {
	var err error
	lUser := list.New()
	Lusr.LogUser = &userError{"init error", nil, fl}
	rows, err := Db.Query("SELECT id_user, mail, ST_AsText(location_user) FROM \"user\";")
	if err != nil {
		fl.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var idUser int64
		var mailq string
		var pos string
		err = rows.Scan(&idUser, &mailq, &pos)
		if err != nil {
			fl.Println(err)
		}
		lUser.PushBack(&User{Id: idUser, Mail: mailq, Followed: list.New(), Stats: Lusr.GetStatsByUser(idUser, Db), HistoricReq: list.New(), Possessed: list.New(), Coord: GetCoord(pos)})
	}
	Lusr.Ulist.Init()
	Lusr.Ulist.PushFrontList(lUser)
	return nil
}
