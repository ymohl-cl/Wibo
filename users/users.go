package users

import (
	"Wibo/protocol"
	"bytes"
	"container/list"
	"database/sql"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"log"
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
	Psd         []byte
}

type All_users struct {
	Ulist      *list.List
	GlobalStat *StatsUser /* Stats globale a tous les utilisateur de WIbo */
	NbrUsers   int64
	Logger     *log.Logger
}

func (User *User) MagnetisValid() bool {
	if time.Since(User.Magnet) > (1 * time.Hour) {
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
	for euser := lst.Front(); euser != nil; euser = euser.Next() {
		user := euser.Value.(*list.Element).Value.(*User)
		if strings.Compare(user.Mail, email) == 0 {
			return euser.Value.(*list.Element)
		}
	}
	return nil
}

func FoundUserOnListLvl1(lst *list.List, email string) *list.Element {
	for euser := lst.Front(); euser != nil; euser = euser.Next() {
		user := euser.Value.(*User)
		if strings.Compare(user.Mail, email) == 0 {
			return euser
		}
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
		user = ulist.CheckPasswordUser(user, req.Spec.(protocol.Log).Pswd, Db)
	}
	return user
}

/******************************************************************************/
/********************************* MERGE JAIME ********************************/
/******************************************************************************/

func CheckValidMail(email string) bool {
	flag_bool := valid.IsEmail(email)
	return flag_bool
}

func (ulist *All_users) CheckPasswordUser(user *list.Element, pass []byte, Db *sql.DB) *list.Element {
	var flag int

	err := Db.QueryRow("SELECT login($1, $2);", user.Value.(*User).Id, pass).Scan(&flag)
	if err != nil {
		ulist.Logger.Println("DbQuery return err:", err)
	}
	if flag == 1 {
		return user
	}
	return nil
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

func (Lst_users *All_users) Add_new_user(new_user *User, Db *sql.DB, Pass []byte) (bool, error) {
	if len(new_user.Mail) > 0 {
		if valid.IsEmail(new_user.Mail) != true {
			return false, nil
		}
	}
	err := Db.QueryRow("SELECT  setsuserdata2($1, $2, $3, $4, $5, $6, $7, $8);", 1, "user_particulier", new_user.Coord.Lat, new_user.Coord.Lon, new_user.Stats.CreationDate, new_user.Log, new_user.Mail, Pass).Scan(&new_user.Id)
	if err != nil {
		Lst_users.Logger.Println("Add_new_user: ", err)
		return false, err
	}
	return true, nil
}

/*
Insert user default
	insert a user with default data
*/
func (Lst_users *All_users) AddNewDefaultUser(Db *sql.DB, req *protocol.Request) *list.Element {
	var err error
	t := bytes.NewBufferString("1")
	bpass := t.Bytes()
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
	err = Db.QueryRow("SELECT setsdefaultuserdata($1, $2, $3, $4, $5);", tmpUser.Coord.Lat, tmpUser.Coord.Lon, tmpUser.Log, tmpUser.Stats.CreationDate, bpass).Scan(&tmpUser.Id)
	if err != nil {
		Lst_users.Logger.Println(err)
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
	rows, err := Db.Query("SELECT id_user, mail FROM \"user\" WHERE id_user=$1;", idUser)
	if err != nil {
		LstU.Logger.Println("Error Query on SelectUser: ", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var idUser int64
		var mailq string
		err = rows.Scan(&idUser, &mailq)
		if err != nil {
			LstU.Logger.Println("Error Scan on SelectUser: ", err)
			return nil
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
		Lusr.Logger.Println("GetDevicesByIdUser query error: ", err)
		return lDevice
	}
	defer stm.Close()
	rows, err := stm.Query(idUser)
	if err != nil {
		Lusr.Logger.Println("GetDevicesByIdUser query error: ", err)
		return lDevice
	}
	defer rows.Close()
	for rows.Next() {
		var idDevice string
		err = rows.Scan(&idDevice)
		if err != nil {
			Lusr.Logger.Println("GetDevicesByIdUser scan error: ", err)
		} else {
			lDevice.PushBack(idDevice)
		}
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

/*func (Lusr *All_users) setBackendLog() {
	backendUser := logging.SetForeingLogBackend(Lusr.Logger)
	backendFormatter := logging.NewBackendFormatter(backendUser, format)
	backend1Leveled := logging.AddModuleLevel(backendUser)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backendFormatter)
}*/

/**
* Get_users
* Query the user table join device and create new *listList Pointer
**/
func (Lusr *All_users) Get_users(Db *sql.DB) (err error) {
	lUser := list.New()
	rows, err := Db.Query("SELECT id_user, mail, ST_AsText(location_user) FROM \"user\";")
	if err != nil {
		Lusr.Logger.Println("GetUser query error: ", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var idUser int64
		var mailq, pos string
		err = rows.Scan(&idUser, &mailq, &pos)
		if err != nil {
			Lusr.Logger.Println("Get_users scan error: ", err)
			return err
		}
		lUser.PushBack(&User{Id: idUser, Mail: mailq, Followed: list.New(), Stats: Lusr.GetStatsByUser(idUser, Db), HistoricReq: list.New(), Possessed: list.New(), Coord: GetCoord(pos)})
	}
	Lusr.Ulist = list.New()
	Lusr.Ulist.PushFrontList(lUser)
	Lusr.Logger.Println("Get number user: ", Lusr.Ulist.Len())
	return nil
}
