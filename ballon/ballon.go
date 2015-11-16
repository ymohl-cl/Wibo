package ballon

import (
	"Wibo/db"
	"Wibo/owm"
	"Wibo/protocol"
	"Wibo/users"
	"container/list"
	"database/sql"
	"fmt"
	af "github.com/spf13/afero"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	NBRCHECKPOINTLIST = 3
)

var SubDir = "/tmp/MeMapUsers"
var fileUsers = "mmu.txt"
var Fss = []af.Fs{&af.MemMapFs{}, &af.OsFs{}}

/* Type is message type. Only type 1 is use now and described a text */
/*
** Structure message qui contient le contenu de la string sa taille et son type
** Pour le moment seul le type 1 est pris en compte pour le texte.
** Le premier message de la liste est le message de creation du ballon.
 */
type Message struct {
	Id        int32
	Size      int32
	Idcountry int32
	Idcity    int32
	Content   string
	Type      int32
}

/* Longitude and latitude */
type Coordinate struct {
	Lon float64
	Lat float64
}

// MagnetFlag == 1 if take with magnet
type Checkpoint struct {
	Coord      Coordinate
	Date       time.Time
	MagnetFlag int16
}

type Wind struct {
	Speed   float64
	Degress float64
}

type StatsBall struct {
	CreationDate time.Time   /* Date de creation, set par le serveur */
	CoordCreated *Coordinate /* Lieu de creation */
	NbrKm        float64     /* Nombre de kilometre parcourus */
	NbrFollow    int64       /* Nombre de personne qui follow le ballon */
	NbrCatch     int64       /* Nombre de fois ou le ballon a ete attrappe. */
	NbrMagnet    int64       /* Nombre de fois ou le ballon a ete aimante */
}

type Ball struct {
	sync.RWMutex
	Id_ball     int64         /* Id de la base de donnee, defini par le serveur */
	Title       string        /* Titre du ballon */
	Coord       *list.Element /* Interface coordonnee, stocke les coordonnees */
	Scoord      *list.Element /** Last itinerary save **/
	Itinerary   *list.List    /* List des itineraire non enregistre par le serveur */
	Edited      bool          /** Edited, Flag de modification (ne tiens pas compte des changement dans Coord et Checkpoints) */
	Wind        Wind          /* Interface Wind, stocke les donnees des vents */
	Messages    *list.List    /* Liste d'interface Message: Contenu du message, l'id 0 est le message de creation du ballon */
	Date        time.Time     /* Creation date, set par le serveur */
	Checkpoints *list.List    /** Liste interface checkpoints */
	Possessed   *list.Element /* Value: (*users.User), user qui possede le ballon */
	Followers   *list.List    /* List d'insterface *list.Element.Value.(*users.User), Constitue la liste des utilisateurs qui suivents le ballon */
	Creator     *list.Element /* Value: (*users.User), user qui a creer le ballon */
	Stats       *StatsBall    /* Interface Stats: Statistiques de la vie du ballon */
	FlagC       bool          // Flag de creation to insert if true or update for false
}

type All_ball struct {
	sync.RWMutex
	Blist  *list.List /* Value: *Ball */
	Id_max int64      /* Set by bdd and incremented by server */
	Logger *log.Logger
	Ftmp   *af.InMemoryFile
}

/*
** Source: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
 */
func hsin(theta float64) (result float64) {
	result = math.Pow(math.Sin(theta/2), 2)
	return
}

/*
** Source: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
 */
func (ball *Ball) GetDistance(lon_user float64, lat_user float64) float64 {
	var lat1, lat2, lon1, lon2, rayon float64
	lat1 = lat_user * math.Pi / 180
	lon1 = lon_user * math.Pi / 180
	lat2 = ball.Coord.Value.(Checkpoint).Coord.Lat
	lon2 = ball.Coord.Value.(Checkpoint).Coord.Lon
	lat2 = ball.Coord.Value.(Checkpoint).Coord.Lat * math.Pi / 180
	lon2 = ball.Coord.Value.(Checkpoint).Coord.Lon * math.Pi / 180
	rayon = 6378137

	hvsin := hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*hsin(lon2-lon1)
	return (2 * rayon * math.Asin(math.Sqrt(hvsin))) / 1000
}

func (ball *Ball) Print_list_checkpoints() {
	echeck := ball.Checkpoints.Front()

	for echeck != nil {
		fmt.Println(echeck.Value.(Checkpoint))
		echeck = echeck.Next()
	}
}

func (ball *Ball) Check_userfollower(user *list.Element) bool {
	euser := ball.Followers.Front()

	for euser != nil && euser.Value.(*list.Element).Value.(*users.User).Id != user.Value.(*users.User).Id {
		euser = euser.Next()
	}
	if euser != nil {
		return true
	}
	return false
}

func (ball *Ball) Check_userCreated(user *list.Element) bool {
	if ball.Creator == user {
		return true
	} else {
		return false
	}
}

func (ball *Ball) Check_userPossessed(user *list.Element) bool {
	if ball.Possessed == user {
		return true
	}
	return false
}

func (ball *Ball) Check_nearbycoord(request *list.Element) bool {
	rlon := request.Value.(*protocol.Request).Coord.Lon
	rlat := request.Value.(*protocol.Request).Coord.Lat

	if ball.Coord != nil && ball.Possessed == nil {
		if ball.GetDistance(rlon, rlat) < 1.0 {
			return true
		}
	}
	return false
}

func (balls *All_ball) Get_ballbyid(id int64) (eball *list.Element) {
	eball = balls.Blist.Front()

	for eball != nil && eball.Value.(*Ball).Id_ball != id {
		eball = eball.Next()
	}
	return eball
}

func EballAlreadyExist(blst *list.List, eball *list.Element) bool {
	ball := eball.Value.(*Ball)
	for e := blst.Front(); e != nil; e = e.Next() {
		if e.Value.(*list.Element).Value.(*Ball).Id_ball == ball.Id_ball {
			return true
		}
	}
	return false
}

func (balls *All_ball) Get_ballbyid_tomagnet(tab [3]int64, User *list.Element) *list.List {
	list_tmp := list.New()
	flag := false

	for i := 0; i < 3; i++ {
		eball := balls.Blist.Front()
		for eball != nil && eball.Value.(*Ball).Id_ball != tab[i] {
			eball = eball.Next()
		}
		for eball != nil && (eball.Value.(*Ball).Possessed != nil || EballAlreadyExist(list_tmp, eball) == true || eball.Value.(*Ball).Check_userCreated(User) == true || eball.Value.(*Ball).Check_userfollower(User) == true) {
			eball = eball.Next()
			if eball == nil && flag == false {
				flag = true
				eball = balls.Blist.Front()
			}
		}
		if eball != nil {
			list_tmp.PushBack(eball)
		}
		flag = false
	}
	return list_tmp
}

/**
** Implementation a grande echelle
**/

func (Ball *Ball) GetCheckpoint(station owm.Weather_data) Checkpoint {
	r_world := 6378137.0
	var tmp_coord Coordinate
	var calc_coord Coordinate
	var checkpoint Coordinate

	speed := station.Wind.Speed * 300.00
	dir := station.Wind.Degress*10 + 180
	if dir >= 360 {
		dir -= 360
	}
	dir = dir * (math.Pi / 180.0)
	if Ball.Checkpoints.Len() == 0 {
		checkpoint.Lon = Ball.Coord.Value.(Checkpoint).Coord.Lon
		checkpoint.Lat = Ball.Coord.Value.(Checkpoint).Coord.Lat
	} else {
		checkpoint = ((Ball.Checkpoints.Back()).Value.(Checkpoint)).Coord
	}
	tmp_coord.Lon = checkpoint.Lon * (math.Pi / 180.0)
	tmp_coord.Lat = checkpoint.Lat * (math.Pi / 180.0)
	calc_coord.Lat = math.Asin(math.Sin(tmp_coord.Lat)*math.Cos(speed/r_world) + math.Cos(tmp_coord.Lat)*math.Sin(speed/r_world)*math.Cos(dir))
	calc_coord.Lon = tmp_coord.Lon + math.Atan2(math.Sin(dir)*math.Sin(speed/r_world)*math.Cos(tmp_coord.Lat), math.Cos(speed/r_world)-math.Sin(tmp_coord.Lat)*math.Sin(calc_coord.Lat))
	calc_coord.Lat = 180 * calc_coord.Lat / math.Pi
	calc_coord.Lon = 180 * calc_coord.Lon / math.Pi
	return Checkpoint{calc_coord, time.Now(), 0}
}

func (Ball *Ball) CreateCheckpoint(Lst_wd *owm.All_data) error {
	var station owm.Weather_data

	Lon := Ball.Coord.Value.(Checkpoint).Coord.Lon
	Lat := Ball.Coord.Value.(Checkpoint).Coord.Lat
	station = Lst_wd.GetNearest(Lon, Lat)

	for i := 0; i < NBRCHECKPOINTLIST; i++ {
		Ball.Checkpoints.PushBack(Ball.GetCheckpoint(station))
	}
	Ball.Wind.Speed = station.Wind.Speed
	Ball.Wind.Degress = station.Wind.Degress
	Ball.Coord = Ball.Checkpoints.Front()
	Ball.Checkpoints.Remove(Ball.Coord)
	return nil
}

func (Lst_ball *All_ball) Create_checkpoint(Lst_wd *owm.All_data) error {
	for eb := Lst_ball.Blist.Front(); eb != nil; eb = eb.Next() {
		ball := eb.Value.(*Ball)
		if ball.Possessed == nil {
			ball.Lock()
			ball.CreateCheckpoint(Lst_wd)
			ball.Unlock()
		}
	}
	return nil
}

func (Ball *Ball) GetTimeTrueCoord() {
	lst := list.New()
	var check Checkpoint

	check.Coord = Ball.Coord.Value.(Checkpoint).Coord
	check.Date = time.Now()
	check.MagnetFlag = Ball.Coord.Value.(Checkpoint).MagnetFlag
	Ball.Coord = lst.PushFront(check)
}

func (Ball *Ball) InitCoord(Lon float64, Lat float64, Magnet int16, Wd *owm.All_data, CrtCK bool) {
	var check Checkpoint
	lst := list.New()

	Ball.Lock()
	check.Coord.Lon = Lon
	check.Coord.Lat = Lat
	check.Date = time.Now()
	check.MagnetFlag = Magnet
	Ball.Coord = lst.PushBack(check)
	Ball.Scoord = Ball.Coord
	Ball.Itinerary.PushBack(Ball.Coord.Value.(Checkpoint))
	Ball.Coord = lst.PushBack(check)
	Ball.Scoord = Ball.Coord
	Ball.Itinerary.PushBack(Ball.Coord.Value.(Checkpoint))
	if CrtCK == true {
		Ball.CreateCheckpoint(Wd)
	}
	Ball.Unlock()
}

func (Lst_ball *All_ball) Move_ball(Lst_wd *owm.All_data) (er error) {
	var coord Coordinate

	fmt.Println("!!!! MOVE BAL !!!!")
	for eb := Lst_ball.Blist.Front(); eb != nil; eb = eb.Next() {
		ball := eb.Value.(*Ball)
		if ball.Possessed == nil {
			ball.Lock()
			fmt.Println("Ball Title: ", ball.Title)
			fmt.Println("Coord de ball avant move: ", ball.Coord.Value.(Checkpoint))
			if ball.Checkpoints.Len() == 0 {
				fmt.Println("Checkpoint empty")
				coord = ball.Coord.Value.(Checkpoint).Coord
				ball.CreateCheckpoint(Lst_wd)
				fmt.Println("CreateCheckpoint")
				ball.Stats.NbrKm += ball.GetDistance(coord.Lon, coord.Lat)
				fmt.Println("Get distance ok")
			} else {
				fmt.Println("Checkpoint no empty")
				e := ball.Checkpoints.Front()
				coord = e.Value.(Checkpoint).Coord
				ball.Stats.NbrKm += ball.GetDistance(coord.Lon, coord.Lat)
				ball.Coord = e
				ball.Checkpoints.Remove(ball.Coord)
				ball.GetTimeTrueCoord()
			}
			Lon := ball.Scoord.Value.(Checkpoint).Coord.Lon
			Lat := ball.Scoord.Value.(Checkpoint).Coord.Lat
			if ball.Itinerary.Len() == 0 || ball.GetDistance(Lon, Lat) > 1.0 {
				fmt.Println("Add on Itinerary")
				ball.Itinerary.PushBack(ball.Coord.Value.(Checkpoint))
				ball.Scoord = ball.Coord
			}
			fmt.Println("Coord de ball apres move: ", ball.Coord.Value.(Checkpoint))
			ball.Unlock()
		}
	}
	return nil
}

/**
** Fin de l'implementation a grande echelle
**/
/**
** Cette section est implemente pour la beta uniquement.
**/
/* Apply the function Get_checkpointlist all ballons */
/*func (Lst_ball *All_ball) Create_checkpointBeta(Lst_wd *owm.All_data) error {
	var station owm.Weather_data

	station = Lst_wd.Get_Paris()
	Lst_ball.Lock()
	defer Lst_ball.Unlock()
	eball := Lst_ball.Blist.Front()
	for eball != nil {
		eball.Value.(*Ball).Get_checkpointList(station)
		eball = eball.Next()
	}
	return nil
}*/
/**
** Fin de la section Beta
**/

/* Add_new_ballon to list */
func (Lst_ball *All_ball) Add_new_ballon(new_ball Ball) {
	Lst_ball.Blist.PushBack(new_ball)
	return
}

func Print_all_message(lst *list.List) {
	emess := lst.Front()

	for emess != nil {
		mess := emess.Value.(Message)
		fmt.Println("Message ...")
		fmt.Println(mess.Id)
		fmt.Println(mess.Size)
		fmt.Println(mess.Content)
		emess = emess.Next()
	}
}

func Print_all_checkpoints(check *list.List) {
	echeck := check.Front()

	for echeck != nil {
		tcheck := echeck.Value.(Checkpoint)
		fmt.Println("Checkpoint ...")
		fmt.Println(tcheck.Coord.Lon)
		fmt.Println(tcheck.Coord.Lat)
		fmt.Println(tcheck.Date)
		echeck = echeck.Next()
	}
}

func Print_users_follower(ulist *list.List) {
	euser := ulist.Front()

	for euser != nil {
		user := euser.Value.(*list.Element).Value.(*users.User)
		fmt.Println("User ...")
		fmt.Println(user.Mail)
		euser = euser.Next()
	}
}

/* Print_all_ball print for debug. */
func (Lst_ball *All_ball) Print_all_balls() {
	eball := Lst_ball.Blist.Front()

	for eball != nil {
		ball := eball.Value.(*Ball)
		fmt.Println("!!!! Print BALL !!!!")
		fmt.Println(ball.Id_ball)
		fmt.Println(ball.Title)
		fmt.Println(ball.Coord)
		fmt.Println(ball.Wind)
		fmt.Println("!!!! MESSAGE !!!!")
		Print_all_message(ball.Messages)
		fmt.Println(ball.Date)
		fmt.Println("!!!! Checkpoints !!!!")
		Print_all_checkpoints(ball.Checkpoints)
		fmt.Println("!!!! User possessed !!!!")
		fmt.Println(ball.Possessed)
		fmt.Println("!!!! Users follower !!!!")
		Print_users_follower(ball.Followers)
		fmt.Println("!!!! User creator !!!!")
		fmt.Println(ball.Creator)
		eball = eball.Next()
	}
	return
}

/******************************************************************************/
/******************************** MERGE JAIME *********************************/
/******************************************************************************/
/*
insert checkpoints(
       cdate date,
       latitudec double precision,
       longitudec double precision,
       idcont integer,
       magnet boolean)
*/
/* Set on itinary list, creation point of ball */
func (ball *Ball) SetCreationCoordOnItinerary(Db *sql.DB) {
	fmt.Println("COUCOU")
	var Idb int64
	coord := ball.Stats.CoordCreated
	fmt.Println("Coord que je veux voir: ", coord)
	row, err := Db.Query("SELECT id from container WHERE ianix = $1", ball.Id_ball)
	if err != nil {
		fmt.Println("Erreur HAHA!", err)
		log.Print(err)
	}
	defer row.Close()
	if row.Next() != false {
		row.Scan(&Idb)
		trow, err := Db.Query("SELECT insertcheckpoints($1, $2, $3, $4, $5)", ball.Stats.CreationDate, coord.Lon, coord.Lat, Idb, 0)
		if err != nil {
			fmt.Println("Erreur HAHA2!", err)
			fmt.Println(err)
		}
		trow.Close()
	} else {
		fmt.Println("Erreur HAHA3!", err)
		fmt.Println(err)
	}
}

func (Lb *All_ball) SetItinerary(Db *sql.DB, b *list.Element) {
	//	for b := Lb.Blist.Front(); b != nil; b = b.Next() {
	var Idb int64
	row, err := Db.Query("SELECT id from container WHERE ianix = $1", b.Value.(*Ball).Id_ball)
	if err != nil {
		log.Print(err)
	}
	defer row.Close()
	if row.Next() != false {
		row.Scan(&Idb)
		for i := b.Value.(*Ball).Itinerary.Front(); i != nil; i = i.Next() {
			trow, err := Db.Query("SELECT insertcheckpoints($1, $2, $3, $4, $5)", i.Value.(Checkpoint).Date, i.Value.(Checkpoint).Coord.Lon, i.Value.(Checkpoint).Coord.Lat, Idb, i.Value.(Checkpoint).MagnetFlag)
			if err != nil {
				fmt.Println(err)
			}
			trow.Close()
		}
		b.Value.(*Ball).Itinerary.Init()
	} else {
		fmt.Println(err)
	}
	b.Value.(*Ball).Itinerary = list.New()
}

func (Ball *Ball) GetItinerary(Db *sql.DB) (int32, *list.List) {
	var err error
	Itinerary := list.New()
	//	Ball.Itinerary = list.New()
	var idB int64
	Db.QueryRow("SELECT id FROM container WHERE ianix=$1;", Ball.Id_ball).Scan(&idB)
	rows, err := Db.Query("SELECT getitenirarybycontainerid($1);", idB)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	//	if rows.Next() != false {
	for rows.Next() {
		var table string
		var ism int16
		rows.Scan(&table)
		result := strings.Split(table, ",")
		ism = 0
		tdate := GetDateFormat(result[0])
		if strings.ContainsRune(result[1], 't') == true {
			ism = 1
		}
		tempCoord := getExtraInfo(result[2], tdate, ism)
		checkp := new(Checkpoint)
		checkp.Coord.Lon = tempCoord.Front().Value.(Checkpoint).Coord.Lon
		checkp.Coord.Lat = tempCoord.Front().Value.(Checkpoint).Coord.Lat
		checkp.Date = tempCoord.Front().Value.(Checkpoint).Date
		checkp.MagnetFlag = ism
		Itinerary.PushBack(checkp)
	}
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}
	// Delete first elem. First elem is a creation point
	fmt.Println("elem ?", err)
	elem := Itinerary.Front()
	if elem != nil {
		fmt.Println("Yes")
		Itinerary.Remove(elem)
	} else {
		fmt.Println("No")
	}
	return int32(Itinerary.Len()), Itinerary
}

func getIdMessageMax(idBall int64, base *db.Env) (int32, error) {
	var IdMax int32
	err := base.Transact(base.Db, func(tx *sql.Tx) error {
		var err error
		stm, err := tx.Prepare("select id from message where id = (select max(id) from message) and containerid = $1;")
		if err != nil {
			return err
		}
		rs, err := stm.Query(idBall)
		if err != nil {
			return err
		}
		defer stm.Close()
		if rs.Next() != false {
			rs.Scan(&IdMax)
		}
		return err
	})
	return IdMax, err
}

func getIdBallMax(base *db.Env) int64 {
	var IdMax int64
	IdMax = 0
	rs, err := base.Db.Query("SELECT ianix FROM container ORDER BY ianix DESC LIMIT 1;")
	if err != nil {
		return IdMax
	}
	defer rs.Close()
	if rs.Next() != false {
		rs.Scan(&IdMax)
	}
	return IdMax
}

/**
* CheckErr
* Verify err value to stop execution by panic
**/
func (Lst_ball *All_ball) checkErr(err error) {
	if err != nil {
		Lst_ball.Logger.Printf("Error: %s", err)
	}
}

func (Lb *All_ball) GetFollowers(idBall int64, Db *sql.DB, Ulist *list.List) (*list.List, error) {
	lstFollow := list.New()
	var err error
	rows, err := Db.Query("SELECT id_user FROM \"user\" AS userWibo LEFT OUTER JOIN followed ON (followed.iduser = userWibo.id_user)  WHERE followed.container_id = $1;", idBall)
	if err != nil {
		return lstFollow, err
	}
	defer rows.Close()
	for rows.Next() {
		var idFollower int64
		rows.Scan(&idFollower)
		for u := Ulist.Front(); u != nil; u = u.Next() {
			if idFollower == u.Value.(*users.User).Id {
				lstFollow.PushBack(u)
			}
		}
	}
	return lstFollow, err
}

func GetCurrentUserBall(LUser *list.List, idBall int64, Db *sql.DB) (*list.Element, error) {
	stm, err := Db.Prepare("SELECT idcurrentuser  FROM container WHERE id=($1)")
	if err != nil {
		return nil, err
	}
	rows, err := stm.Query(idBall)
	if err != nil {
		return nil, err
	}

	if rows.Next() != false {
		var idPossesed int64
		err = rows.Scan(&idPossesed)
		if err != nil {
			return nil, err
		}
		i := 0
		for e := LUser.Front(); e != nil; e = e.Next() {
			if e.Value.(*users.User).Id == idPossesed {
				return e, err
			}
			i++
		}
	}
	return nil, err
}

func GetWhomGotBall(idBall int64, LstU *list.List, Db *sql.DB) (*list.Element, error) {
	p, er := GetCurrentUserBall(LstU, idBall, Db)
	return p, er
}

func getExtraInfo(position string, date time.Time, magnet int16) *list.List {

	f := func(c rune) bool {
		return c == '(' || c == '(' || c == ')' || c == '"' ||
			c == 'P' || c == 'O' || c == 'I' || c == 'N' ||
			c == 'T'
	}
	// Separate into fields with func.
	fmt.Println("Warning, field: ", position)
	fields := strings.FieldsFunc(position, f)
	// Separate into cordinates  with Fields.
	point := strings.Fields(fields[0])
	long, err := strconv.ParseFloat(point[0], 15)
	if err != nil {
		fmt.Println(err)
	}
	var lat float64
	lat, _ = strconv.ParseFloat(point[1], 15)
	lc := list.New()
	ch := Coordinate{Lon: long, Lat: lat}

	//	newdate := GetDateFormat(date)
	//	tmagnet, _ := strconv.Atoi(magnet)
	lc.PushBack(Checkpoint{Coord: ch, Date: date, MagnetFlag: magnet})
	return lc
}

/**
* GetListBallsByUser
* getContainersByUserId is a native psql function with
* RETURNS TABLE(
			idballon integer,
			titlename varchar(255),
			idtype integer, direction numeric,
			speedcont integer,
			creationdate date,
			deviceid integer,
 CREATE OR REPLACE FUNCTION public.getcontainersbyuserid(iduser integer)
  RETURNS TABLE(idballon integer, titlename character varying, idtype integer, direction float, speedcont float, creationdate timestamp without time zone, locationcont text)
  LANGUAGE plpgsql
 AS $function$  BEGIN RETURN QUERY SELECT container.id AS contIndex, container.titlename AS    TitleName, container.id_type_c AS TypeCode, container.direction AS contDirection, container.speed AS contSpeed, date(container.creationdate) + interval '1 hour',   ST_AsText(container.location_ct)  FROM container  WHERE idcreator = iduser;  END $function$
			locationcont text)
*/

func (Lb *All_ball) GetListBallsByUser(userE *list.Element, base *db.Env, Ulist *list.List) (lBallon *list.List, err error) {
	lBallon = list.New()
	err = nil

	err = base.Transact(base.Db, func(tx *sql.Tx) error {
		var errT error
		stm, errT := tx.Prepare("SELECT getcontainersbyuserid($1);")
		if errT != nil {
			Lb.Logger.Println(err)
			//return err
		}
		rows, err := stm.Query(userE.Value.(*users.User).Id)
		if err != nil {
			Lb.Logger.Println(err)
			//	return err
		}
		if rows.Next() != false {
			for rows.Next() {
				var infoCont string
				err = rows.Scan(&infoCont)
				if err != nil {
					Lb.Logger.Println(err)
					//		return err
				}
				result := strings.Split(infoCont, ",")
				fmt.Printf("%T | %v \n", infoCont, infoCont)
				idBall := GetIdBall(result[0])
				magnet, _ := strconv.Atoi(result[8])
				tempCord := getExtraInfo(result[6], GetDateFormat(result[9]), int16(magnet))
				lstIt := list.New()
				//lstIt.PushFront(tempCord.Front().Value.(Checkpoint))
				possessed, _ := GetWhomGotBall(idBall, Ulist, base.Db)
				if possessed == nil {
					possessed = userE
					Lb.Logger.Println("Possesed is null")
					//		return er
				}
				tmpBall := Lb.Get_ballbyid(GetIdBall(result[7]))
				fmt.Printf("%v and %v \n", tmpBall, result[7])
				if tmpBall != nil {
					// Do Nothing
				} else {
					lstMess, err := Lb.GetMessagesBall(idBall, base.Db)
					if err != nil {
						Lb.Logger.Println(err)
						//			return err
					}
					lstFols, err := Lb.GetFollowers(idBall, base.Db, Ulist)
					if err != nil {
						Lb.Logger.Println(err)
						//			return err
					}
					lBallon.PushBack(
						&Ball{
							Id_ball:     GetIdBall(result[7]),
							Title:       result[1],
							Date:        GetDateFormat(result[5]),
							Checkpoints: nil,
							Itinerary:   lstIt,
							Scoord:      tempCord.Front(),
							Coord:       tempCord.Front(),
							Wind:        GetWin(result[3], result[4]),
							Messages:    lstMess,
							Followers:   lstFols,
							Possessed:   possessed,
							Stats:       Lb.GetStatsBallon(int64(idBall), base.Db),
							Creator:     userE})
				}
			}
		}
		stm.Close()
		return err
	})
	if err != nil {
		Lb.Logger.Println(err)
		return nil, err
	}
	return lBallon, nil

}

func GetDateFormat(qdate string) (fdate time.Time) {
	f := func(c rune) bool {
		return c == '"'
	}
	fields := strings.FieldsFunc(qdate, f)
	for _, value := range fields {
		qdate = string(value)
	}
	fdate, _ = time.Parse("2006-01-02 15:04:05", qdate)
	return fdate
}

/**
* GetIdBall
* Parse and []string with the id Value
* Get and array parsed and convert this value
* return the id
 */
func GetIdBall(idB string) int64 {
	// Return true if 'value' char.
	f := func(c rune) bool {
		return c == '(' || c == '(' || c == ')' || c == '"'
	}
	// Separate into fields with func.
	fields := strings.FieldsFunc(idB, f)
	// Separate into cordinates  with Fields.
	ids := strings.Fields(fields[0])
	id, _ := strconv.Atoi(ids[0])
	return int64(id)
}

/**
* GetCord
* Parse query function POINT(longitude latitude)
* Get the Values between the parenthesis
* convert them to float
* create and new Cooridantes element and return it
**/
func GetCord(position string) (coord *Coordinate) {
	var err error
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
	coord.Lon, err = strconv.ParseFloat(point[0], 15)
	if err != nil {
		fmt.Println(err)
	}
	coord.Lat, _ = strconv.ParseFloat(point[1], 15)
	return coord
}

/**
* GetWin
* Take to strings to Parse their value to Float
* return a new Instace of Wind with float values
**/

func GetWin(speed string, direction string) Wind {
	sf, _ := strconv.ParseFloat(speed, 10)
	df, _ := strconv.ParseFloat(direction, 10)
	return (Wind{Speed: sf, Degress: df})
}

/**
* GetMessageBall
* Query the message who concern an idContainer by timestamp
* Create a new list of message by id Container
* Create an element of list
* Push back this element in a new  list of message
* return the list
**/

func (Lball *All_ball) GetMessagesBall(idBall int64, Db *sql.DB) (*list.List, error) {
	Mlist := list.New()

	stm, err := Db.Prepare("SELECT id AS containerId, content, id_type_m  FROM message WHERE containerid=($1) ORDER BY creationdate DESC")
	defer stm.Close()
	rows, err := stm.Query(idBall)
	if err != nil {
		return Mlist, err
	}
	for rows.Next() {
		var idm int32
		var message string
		var idType int32
		err = rows.Scan(&idm, &message, &idType)
		if err != nil {
			return Mlist, err
		}
		Mlist.PushBack(&Message{Content: message, Type: idType, Id: idm})
	}
	return Mlist, err
}

func (Lball *All_ball) InsertListBallsFollow(Blist *list.List, Ulist *list.List, base *db.Env) {
	for b := Blist.Front(); b != nil; b = b.Next() {
		for f := b.Value.(*Ball).Followers.Front(); f != nil; f = f.Next() {
			f.Value.(*list.Element).Value.(*users.User).Followed.PushBack(b)
		}
	}

	//	for u := Ulist.Front(); u != nil; u = u.Next() {
	//		for b := Blist.Front(); b != nil; b = b.Next() {
	//			for f := b.Value.(*Ball).Followers.Front(); f != nil; f = f.Next() {
	//				if f.Value.(*users.User).Id == u.Value.(*users.User).Id {
	//					base.Db.QueryRow("INSERT INTO followed(container_id, iduser) VALUES($1, $2);", b.Value.(*Ball).Id_ball, u.Value.(*users.User).Id)
	//					u.Value.(*users.User).Followed.PushBack(b)
	//				}
	//			}
	//		}
	//	}
}

/**
* get all ball from database and associeted
* the creator, possessord and followers.
**/
func (Lb *All_ball) Get_balls(LstU *users.All_users, base *db.Env) (err error) {
	err = nil

	Lb.Id_max = getIdBallMax(base) + 1
	fmt.Println("Id max de ball: ", Lb.Id_max)
	Lb.Ftmp = af.MemFileCreate("testfile")
	for e := LstU.Ulist.Front(); e != nil; e = e.Next() {
		if e.Value.(*users.User) != nil {
			tlst, err := Lb.GetListBallsByUser(e, base, LstU.Ulist)
			if err != nil {
				return err
			}
			Lb.Blist.PushBackList(tlst)
		}
	}
	Lb.InsertListBallsFollow(Lb.Blist, LstU.Ulist, base)
	return err
}
