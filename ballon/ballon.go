package ballon

import (
	"Wibo/db"
	"Wibo/owm"
	"Wibo/protocol"
	"Wibo/users"
	"container/list"
	"database/sql"
	//	"errors"
	"fmt"
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
	Id_ball int64         /* Id de la base de donnee, defini par le serveur */
	Title   string        /* Titre du ballon */
	Coord   *list.Element /* Interface coordonnee, stocke les coordonnees */
	Scoord  *list.Element /** Last itinerary save **/
	//	Idball      int64         /** Idball -- deprecated a supprimer */
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
}

type All_ball struct {
	sync.RWMutex
	Blist  *list.List /* Value: *Ball */
	Id_max int64      /* Set by bdd and incremented by server */
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
	lat2 = ball.Coord.Value.(Checkpoint).Coord.Lat * math.Pi / 180
	lon2 = ball.Coord.Value.(Checkpoint).Coord.Lon * math.Pi / 180
	rayon = 6378137

	hvsin := hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*hsin(lon2-lon1)
	return (2 * rayon * math.Asin(math.Sqrt(hvsin))) * 1000
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
	if ball.Coord != nil {
		//		coord := ball.Coord.Value.(Checkpoint).Coord
		if ball.GetDistance(rlon, rlat) < 1.0 { // { coord.Lon < rlon+0.01 &&
			//coord.Lon > rlon-0.01 &&
			//coord.Lat < rlat+0.01 &&
			//coord.Lat > rlat-0.01 {
			return true
		}
	}
	return false
}

/* Create list checkpoints for eball with interval time of 5 minutes on 3 hours */
func (eball *Ball) Get_checkpointList(station owm.Weather_data) {
	r_world := 6371000.0
	var tmp_coord Coordinate
	var calc_coord Coordinate
	var checkpoint Coordinate

	speed := station.Wind.Speed * 300.00
	dir := station.Wind.Degress*10 + 180
	if dir >= 360 {
		dir -= 360
	}
	dir = dir * (math.Pi / 180.0)
	checkpoint.Lon = eball.Coord.Value.(Checkpoint).Coord.Lon
	checkpoint.Lat = eball.Coord.Value.(Checkpoint).Coord.Lat
	eball.Checkpoints = eball.Checkpoints.Init()
	for i := 0; i < 35; i++ {
		tmp_coord.Lon = checkpoint.Lon * (math.Pi / 180.0)
		tmp_coord.Lat = checkpoint.Lat * (math.Pi / 180.0)
		calc_coord.Lat = math.Asin(math.Sin(tmp_coord.Lat)*math.Cos(speed/r_world) + math.Cos(tmp_coord.Lat)*math.Sin(speed/r_world)*math.Cos(dir))
		calc_coord.Lon = tmp_coord.Lon + math.Atan2(math.Sin(dir)*math.Sin(speed/r_world)*math.Cos(tmp_coord.Lat), math.Cos(speed/r_world)-math.Sin(tmp_coord.Lat)*math.Sin(calc_coord.Lat))
		calc_coord.Lat = 180 * calc_coord.Lat / math.Pi
		calc_coord.Lon = 180 * calc_coord.Lon / math.Pi
		if calc_coord.Lat < 48.72 {
			checkpoint.Lat = 49.02
		} else if calc_coord.Lat > 49.02 {
			checkpoint.Lat = 48.72
		} else {
			checkpoint.Lat = calc_coord.Lat
		}
		if calc_coord.Lon < 2.10 {
			checkpoint.Lon = 2.60
		} else if calc_coord.Lon > 2.60 {
			checkpoint.Lon = 2.10
		} else {
			checkpoint.Lon = calc_coord.Lon
		}
		eball.Checkpoints.PushBack(Checkpoint{checkpoint, time.Now(), 0})
	}
	eball.Wind.Speed = station.Wind.Speed
	eball.Wind.Degress = station.Wind.Degress
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
		for eball.Value.(*Ball).Possessed != nil || EballAlreadyExist(list_tmp, eball) == true || eball.Value.(*Ball).Check_userCreated(User) == true {
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
	checkpoint.Lon = Ball.Coord.Value.(Checkpoint).Coord.Lon
	checkpoint.Lat = Ball.Coord.Value.(Checkpoint).Coord.Lat
	tmp_coord.Lon = checkpoint.Lon * (math.Pi / 180.0)
	tmp_coord.Lat = checkpoint.Lat * (math.Pi / 180.0)
	calc_coord.Lat = math.Asin(math.Sin(tmp_coord.Lat)*math.Cos(speed/r_world) + math.Cos(tmp_coord.Lat)*math.Sin(speed/r_world)*math.Cos(dir))
	calc_coord.Lon = tmp_coord.Lon + math.Atan2(math.Sin(dir)*math.Sin(speed/r_world)*math.Cos(tmp_coord.Lat), math.Cos(speed/r_world)-math.Sin(tmp_coord.Lat)*math.Sin(calc_coord.Lat))
	calc_coord.Lat = 180 * calc_coord.Lat / math.Pi
	calc_coord.Lon = 180 * calc_coord.Lon / math.Pi
	return Checkpoint{checkpoint, time.Now(), 0}
}

func (Ball *Ball) CreateCheckpoint(Lst_wd *owm.All_data) error {
	var station owm.Weather_data

	Lon := Ball.Coord.Value.(Checkpoint).Coord.Lon
	Lat := Ball.Coord.Value.(Checkpoint).Coord.Lat
	station = Lst_wd.GetNearest(Lon, Lat)

	Ball.Lock()
	for i := 0; i < NBRCHECKPOINTLIST; i++ {
		Ball.Checkpoints.PushBack(Ball.GetCheckpoint(station))
	}
	Ball.Wind.Speed = station.Wind.Speed
	Ball.Wind.Degress = station.Wind.Degress
	Ball.Coord = Ball.Checkpoints.Front()
	Ball.Checkpoints.Remove(Ball.Coord)
	Ball.Unlock()
	return nil
}

func (Lst_ball *All_ball) Create_checkpoint(Lst_wd *owm.All_data) error {
	for eb := Lst_ball.Blist.Front(); eb != nil; eb = eb.Next() {
		ball := eb.Value.(*Ball)
		if ball.Possessed == nil {
			ball.CreateCheckpoint(Lst_wd)
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

	check.Coord.Lon = Lon
	check.Coord.Lat = Lat
	check.Date = time.Now()
	check.MagnetFlag = Magnet
	Ball.Lock()
	Ball.Coord = lst.PushBack(lst)
	Ball.Scoord = Ball.Coord
	Ball.Itinerary.PushBack(Ball.Coord.Value.(Checkpoint))
	Ball.Unlock()
	if CrtCK == true {
		Ball.CreateCheckpoint(Wd)
	}
}

func (Lst_ball *All_ball) Move_ball(Lst_wd *owm.All_data) (er error) {
	fmt.Println("Move Ball")
	for eb := Lst_ball.Blist.Front(); eb != nil; eb = eb.Next() {
		ball := eb.Value.(*Ball)
		if ball.Possessed == nil {
			fmt.Println("Change Coord")
			ball.Lock()
			if ball.Checkpoints.Len() == 0 {
				fmt.Println("Create Checkpoint")
				ball.CreateCheckpoint(Lst_wd)
			} else {
				fmt.Println("Next Checkpoint")
				ball.Coord = ball.Checkpoints.Front()
				ball.Checkpoints.Remove(ball.Coord)
				ball.GetTimeTrueCoord()
			}
			Lon := ball.Scoord.Value.(Checkpoint).Coord.Lon
			Lat := ball.Scoord.Value.(Checkpoint).Coord.Lat
			if ball.Itinerary.Len() == 0 || ball.GetDistance(Lon, Lat) > 1.0 {
				ball.Itinerary.PushBack(ball.Coord.Value.(Checkpoint))
				ball.Scoord = ball.Coord
			}
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
func (Lst_ball *All_ball) Create_checkpointBeta(Lst_wd *owm.All_data) error {
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
}

/**
** Fin de la section Beta
**/

/* Give a next checkpoint ball and removes the previous */
/*
func (Lst_ball *All_ball) Move_ball() (er error) {
	Lst_ball.Lock()
	defer Lst_ball.Unlock()
	elem := Lst_ball.Blist.Front()

	for elem != nil {
		ball := elem.Value.(*Ball)
		if ball.Coord != nil {
			statsCoord := ball.Coord.Value.(Checkpoint).Coord
			ball.Coord = ball.Coord.Next()
			if ball.Coord != nil {
				ball.Checkpoints.Remove(ball.Checkpoints.Front())
			} else {
				ball.Coord = ball.Checkpoints.Front()
				if ball.Coord == nil {
					er = errors.New("next coord not found")
					return er
				}
			}
			ball.Stats.NbrKm += ball.AddStatsDistance(statsCoord.Lon, statsCoord.Lat)
		}
		elem.Value = ball
		elem = elem.Next()
	}
	return nil
}
*/

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
func (Lb *All_ball) SetItinerary(Db *sql.DB) {
	for b := Lb.Blist.Front(); b != nil; b = b.Next() {
		var Idb int64
		row, err := Db.Query("SELECT id from container WHERE ianix = $1", b.Value.(*Ball).Id_ball)
		if err != nil {
			log.Print(err)
		}
		defer row.Close()
		if row.Next() != false {
			row.Scan(&Idb)
			for i := b.Value.(*Ball).Itinerary.Front(); i != nil; i = i.Next() {
				trow, err := Db.Query("SELECT insertcheckpoints($1, $2 $3, $4, 5)", i.Value.(*Checkpoint).Date, i.Value.(*Checkpoint).Coord.Lon, i.Value.(*Checkpoint).Coord.Lat, Idb, i.Value.(*Checkpoint).MagnetFlag)
				if err != nil {
					fmt.Println(err)
				}
				trow.Close()
			}
		} else {
			fmt.Println(err)
		}
		b.Value.(*Ball).Itinerary = list.New()
	}
}

func (Ball *Ball) GetItinerary(Db *sql.DB) (int32, *list.List) {
	var err error
	Ball.Itinerary = list.New()
	rows, err := Db.Query("SELECT date, attractbymagnet, ST_AsText(checkpoints.location_ckp) FROM checkpoints WHERE containerid=$1 ORDER BY date DESC", Ball.Id_ball)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	if rows.Next() != false {
		for rows.Next() {
			var tdate time.Time
			var attm bool
			var point string
			rows.Scan(&tdate, &attm, &point)
			tempCoord := GetCord(point)
			fmt.Println(tdate, attm, point)
			fmt.Printf("GetItinerary %T | %v ", point, point)
			Ball.Itinerary.PushBack(&Checkpoint{Date: tdate, Coord: tempCoord.Front().Value.(Coordinate)})
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return 0, Ball.Itinerary
}

func getIdMessageMax(idBall int64, base *db.Env) int32 {
	var IdMax int32
	err := base.Transact(base.Db, func(tx *sql.Tx) error {
		var err error
		stm, err := tx.Prepare("select id from message where id = (select max(id) from message) and containerid = $1;")
		checkErr(err)
		rs, err := stm.Query(idBall)
		checkErr(err)
		defer stm.Close()
		if rs.Next() != false {
			rs.Scan(&IdMax)
		}
		return err
	})
	checkErr(err)
	return IdMax
}

func getIdBallMax(base *db.Env) int64 {
	var IdMax int64
	IdMax = 0
	rs, err := base.Db.Query("SELECT ianix FROM container ORDER BY ianix DESC LIMIT 1;")
	if err != nil {
		fmt.Println(err)
	}
	defer rs.Close()
	if rs.Next() != false {
		rs.Scan(&IdMax)
	}
	checkErr(err)
	return IdMax
}

/*
 FUNCTION insertContainer(
	$1 idcreatorc integer,
	$2 latitudec integer,
	$3 longitudec integer,
	$5 directionc float,
	$6 speedc float,
	$7 title text,
	$8 idx integer)

	public.insertcontainer(idcreatorc integer,
	$1 latitudec double precision,
	$2 longitudec double precision,
	$3 directionc double precision,
	$4 speedc double precision,
	$5 title text,
	$6 idx integer,
	$7 creation date)
*/

func (Lst_ball *All_ball) InsertBallon(NewBall *Ball, base *db.Env) (executed bool, err error) {
	var IdC int64
	err = base.Transact(base.Db, func(tx *sql.Tx) error {
		stm, err := tx.Prepare("SELECT insertContainer($1, $2, $3, $4, $5, $6 , $7, $8)")
		if err != nil {
			return (err)
		}
		fmt.Printf("insert ballon coordinate lat : type: %T | value: %v\n", NewBall, NewBall.Id_ball)
		fmt.Printf("insert ballon coordinate lat : type: %T | value: %v\n", NewBall.Coord.Value.(Checkpoint).Coord.Lat, NewBall.Coord.Value.(Checkpoint).Coord.Lat)
		fmt.Printf("insert ballon coordinate lon : type: %T | value: %v\n", NewBall.Coord.Value.(Checkpoint).Coord.Lon, NewBall.Coord.Value.(Checkpoint).Coord.Lon)
		fmt.Printf("insert ballon degress : type: %T | value: %v\n", NewBall.Wind.Degress, NewBall.Wind.Degress)
		fmt.Printf("insert ballon Speed : type: %T | value: %v\n", NewBall.Wind.Speed, NewBall.Wind.Speed)
		fmt.Printf("insert ballon title : type: %T | value: %v\n", NewBall.Title, NewBall.Title)
		fmt.Printf("insert ballon idball : type: %T | value: %v\n", NewBall.Id_ball, NewBall.Id_ball)
		fmt.Printf("insert ballon date : type: %T | value: %v\n", NewBall.Date, NewBall.Date)
		err = stm.QueryRow(NewBall.Creator.Value.(*users.User).Id,
			NewBall.Coord.Value.(Checkpoint).Coord.Lat,
			NewBall.Coord.Value.(Checkpoint).Coord.Lon,
			NewBall.Wind.Degress,
			NewBall.Wind.Speed,
			NewBall.Title,
			NewBall.Id_ball,
			NewBall.Date).Scan(&IdC)
		return err
	})
	log.Println(err)
	err = Lst_ball.InsertMessages(NewBall.Messages, IdC, base)
	executed = true
	return executed, err
}

/*
CREATE OR REPLACE FUNCTION public.insertcontainer(idcreatorc integer, latitudec double precision, longitudec double precision, directionc double precision, speedc double precision, title text, idx integer, creation date)
 RETURNS SETOF integer
 LANGUAGE plpgsql
AS $function$  BEGIN RETURN QUERY INSERT INTO container (direction, speed, location_ct, idcreator, titlename, ianix, creationdate) VALUES(directionc, speedc , ST_SetSRID(ST_MakePoint(latitudec, longitudec), 4326), idcreatorc, title, idx, creation) RETURNING id;  END; $function$
\*/
func (Lb *All_ball) Update_balls(ABalls *All_ball, base *db.Env) (er error) {
	i := 0
	fmt.Println("\x1b[31;1m coucou update\x1b[0m")
	fmt.Printf("%v Id Max\n", ABalls.Id_max)
	for e := ABalls.Blist.Front(); e != nil; e = e.Next() {

		if e.Value.(*Ball).Edited == true && e.Value.(*Ball).Id_ball <= ABalls.Id_max {
			e.Value.(*Ball).Lock()

			idBall := e.Value.(*Ball).Id_ball
			idMessageMax := getIdMessageMax(idBall, base)
			j := 0
			for f := e.Value.(*Ball).Messages.Front(); f != nil; f = f.Next() {
				if f.Value.(Message).Id > idMessageMax {
					err := base.Transact(base.Db, func(tx *sql.Tx) error {
						stm, err := tx.Prepare("INSERT INTO message(content, containerid) values($1, (SELECT id from container where ianix = $2))")
						checkErr(err)
						res, err := stm.Exec(f.Value.(Message).Content, idBall)
						var rowsAffect int64
						rowsAffect, err = res.RowsAffected()
						fmt.Println("\x1b[31;1m update %d \x1b[0m", rowsAffect)
						j++
						checkErr(err)
						return err
					})
					checkErr(err)
				}
			}
			e.Value.(*Ball).Unlock()
		} else if e.Value.(*Ball).Id_ball > ABalls.Id_max {
			fmt.Printf("\x1b[31;1m insert ball  %d \x1b[0m\n", e.Value.(*Ball).Id_ball)
			Lb.InsertBallon(e.Value.(*Ball), base)
		}
		i++
	}
	return er
}

func (Lst_ball *All_ball) InsertMessages(messages *list.List, idBall int64, base *db.Env) (err error) {
	i := 0
	for e := messages.Front(); e != nil; e = e.Next() {
		err = base.Transact(base.Db, func(tx *sql.Tx) error {
			stm, err := tx.Prepare("INSERT INTO message(content, containerid) VALUES ($1, $2)")
			checkErr(err)
			_, err = stm.Query(e.Value.(Message).Content, idBall)
			i++
			checkErr(err)
			return err
		})
	}
	return nil
}

/**
* InsertBallonByChamp
* Debug: send an instance of Ball Strut
* Modify the parametres as your needs
**/
func (Lst_ball *All_ball) GetBall(titlename string, Db *sql.DB) *Ball {
	b := new(Ball)
	b.Title = titlename

	u := new(users.User)
	u.Id = 2

	luser := list.New()
	luser.PushBack(u)
	b.Creator = luser.Back()
	b.Id_ball = 5
	b.Wind.Degress = 23.90
	b.Wind.Speed = 222
	b.Messages = Lst_ball.GetMessagesBall(10, Db)
	return (b)
}

/**
* CheckErr
* Verify err value to stop execution by panic
**/
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (Lb *All_ball) GetFollowers(idBall int, Db *sql.DB, Ulist *list.List) *list.List {
	lstFollow := list.New()
	var err error
	rows, err := Db.Query("SELECT id_user FROM \"user\" AS userWibo LEFT OUTER JOIN followed ON (followed.iduser = userWibo.id_user)  WHERE followed.container_id = $1;", idBall)
	checkErr(err)
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
	return lstFollow
}

func GetCurrentUserBall(LUser *list.List, idBall int, Db *sql.DB) *list.Element {
	stm, err := Db.Prepare("SELECT idcurrentuser  FROM container WHERE id=($1)")
	checkErr(err)
	rows, err := stm.Query(idBall)
	checkErr(err)
	for rows.Next() {
		var idPossesed int64
		err = rows.Scan(&idPossesed)
		checkErr(err)
		i := 0
		for e := LUser.Front(); e != nil; e = e.Next() {
			if e.Value.(*users.User).Id == idPossesed {
				return e
			}
			i++
		}
	}
	return nil
}

func GetWhomGotBall(idBall int, LstU *list.List, Db *sql.DB) <-chan *list.Element {
	p := make(chan *list.Element)
	go func() { p <- GetCurrentUserBall(LstU, idBall, Db) }()
	return p
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

func (Lb *All_ball) GetListBallsByUser(userE *list.Element, base *db.Env, Ulist *list.List) *list.List {
	lBallon := list.New()
	var err error
	err = base.Transact(base.Db, func(tx *sql.Tx) error {
		var errT error
		stm, errT := tx.Prepare("SELECT public.getContainersByUserId($1)")
		checkErr(errT)
		fmt.Printf("getlistballsbyuser type %T, value: %v\n", userE.Value.(*users.User).Id, userE.Value.(*users.User).Id)
		rows, err := stm.Query(userE.Value.(*users.User).Id)
		if rows.Next() == false {
			return nil
		}
		switch {
		case err == sql.ErrNoRows:
			log.Printf("No containers.")
		case err != nil:
			log.Print("Error: get containersbyuserid")
			return (err)
		default:
			for rows.Next() {
				var infoCont string
				err = rows.Scan(&infoCont)
				checkErr(err)
				result := strings.Split(infoCont, ",")
				idBall := GetIdBall(result[0])
				tempCord := GetCord(result[7])
				lstIt := list.New()
				lstIt.PushFront(tempCord.Front().Value.(Checkpoint))
				idTmp, _ := strconv.Atoi(result[8])
				possessed := GetWhomGotBall(idBall, Ulist, base.Db)
				tmpBall := Lb.Get_ballbyid(int64(idTmp))
				if tmpBall != nil {

				} else {
					lBallon.PushBack(
						&Ball{
							Title:       result[1],
							Date:        GetDateFormat(result[5]),
							Checkpoints: nil,
							Itinerary:   lstIt,
							Scoord:      tempCord.Front(),
							Coord:       tempCord.Front(),
							Wind:        GetWin(result[3], result[4]),
							Messages:    Lb.GetMessagesBall(idBall, base.Db),
							Followers:   Lb.GetFollowers(idBall, base.Db, Ulist),
							Possessed:   <-possessed,
							Creator:     userE})
				}
				switch {
				case err == sql.ErrNoRows:
					log.Printf("No containers.")
				case err != nil:
					checkErr(err)
					return err
				default:
					fmt.Printf("get containers%v | %v | %v \n\n", result[1], result[3], result[4])
				}
			}
		}
		return err
	})
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No containers.")
	case err != nil:
		log.Print("Here some error")
		log.Print(err)
	default:
		fmt.Printf("get containers 2 %v\n", lBallon)
	}
	return lBallon
}

func GetDateFormat(qdate string) (fdate time.Time) {
	f := func(c rune) bool {
		return c == '"'
	}
	fields := strings.FieldsFunc(qdate, f)
	for _, value := range fields {
		qdate = string(value)
	}
	fdate, err := time.Parse("2006-01-02 15:04:05", qdate)
	checkErr(err)
	return fdate
}

/**
* GetIdBall
* Parse and []string with the id Value
* Get and array parsed and convert this value
* return the id
 */
func GetIdBall(idB string) int {
	// Return true if 'value' char.
	f := func(c rune) bool {
		return c == '(' || c == '(' || c == ')' || c == '"'
	}
	// Separate into fields with func.
	fields := strings.FieldsFunc(idB, f)
	// Separate into cordinates  with Fields.
	ids := strings.Fields(fields[0])
	id, _ := strconv.Atoi(ids[0])
	return (id)
}

/**
* GetCord
* Parse query function POINT(longitude latitude)
* Get the Values between the parenthesis
* convert them to float
* create and new Cooridantes element and return it
**/
func GetCord(position string) *list.List {
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
	lc := list.New()
	lc.PushBack(Coordinate{Lon: long, Lat: lat})
	return lc
}

/**
* GetWin
* Take to strings to Parse their value to Float
* return a new Instace of Wind with float values
**/

func GetWin(speed string, direction string) Wind {
	sf, _ := strconv.ParseFloat(speed, 6)
	df, _ := strconv.ParseFloat(direction, 6)
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

func (Lball *All_ball) GetMessagesBall(idBall int, Db *sql.DB) *list.List {
	Mlist := list.New()
	stm, err := Db.Prepare("SELECT id AS containerId, content, id_type_m  FROM message WHERE containerid=($1) ORDER BY creationdate DESC")
	checkErr(err)
	rows, err := stm.Query(idBall)
	checkErr(err)
	for rows.Next() {
		var idm int32
		var message string
		var idType int32
		err = rows.Scan(&idm, &message, &idType)
		checkErr(err)
		Mlist.PushBack(&Message{Content: message, Type: idType, Id: idm})
	}
	return Mlist
}

func (Lball *All_ball) InsertListBallsFollow(Blist *list.List, Ulist *list.List, base *db.Env) {
	for u := Ulist.Front(); u != nil; u = u.Next() {
		for b := Blist.Front(); b != nil; b = b.Next() {
			for f := b.Value.(*Ball).Followers.Front(); f != nil; f = f.Next() {
				if f.Value.(*users.User).Id == u.Value.(*users.User).Id {
					u.Value.(*users.User).Followed.PushBack(b)
				}
			}
		}
	}
}

/**
* get all ball from database and associeted
* the creator, possessord and followers.
**/
func (Lb *All_ball) Get_balls(LstU *users.All_users, base *db.Env) error {
	Lb.Id_max = getIdBallMax(base)
	for e := LstU.Ulist.Front(); e != nil; e = e.Next() {
		Lb.Blist.PushBackList(Lb.GetListBallsByUser(e, base, LstU.Ulist))
	}
	Lb.InsertListBallsFollow(Lb.Blist, LstU.Ulist, base)
	return nil
}
