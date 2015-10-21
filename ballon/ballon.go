//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  ballon.go                                          :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  By: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  Created: 2015/06/30 16:00:08 by ymohl-cl          #+#    #+#              #
//#  Updated: 2015/06/30 16:00:08 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

package ballon

import (
	"Wibo/db"
	"Wibo/owm"
	"Wibo/protocol"
	"Wibo/users"
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
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
	Id_ball int64         /* Id de la base de donnee, defini par le serveur */
	Title   string        /* Titre du ballon */
	Coord   *list.Element /* Interface coordonnee, stocke les coordonnees */
	//	Idball      int64         /** Idball -- deprecated a supprimer */
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

func (ball *Ball) AddStatsDistance(lon_user float64, lat_user float64) float64 {
	var r float64
	r = 6371000
	var a float64
	var c float64
	var d float64
	var latdiff float64
	var londiff float64

	lonBall := ball.Coord.Value.(Checkpoint).Coord.Lon
	latBall := ball.Coord.Value.(Checkpoint).Coord.Lon
	latdiff = (lat_user - latBall) * (math.Pi / 180)
	londiff = (lon_user - lonBall) * (math.Pi / 180)
	latBall = latBall * (math.Pi / 180)
	lat_user = lat_user * (math.Pi / 180)
	a = math.Sin(latdiff/2)*math.Sin(latdiff/2) + math.Cos(lat_user)*math.Cos(latBall)*math.Sin(londiff/2)*math.Sin(londiff/2)
	c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d = r * c
	return d / 100
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
		coord := ball.Coord.Value.(Checkpoint).Coord
		if coord.Lon < rlon+0.01 &&
			coord.Lon > rlon-0.01 &&
			coord.Lat < rlat+0.01 &&
			coord.Lat > rlat-0.01 {
			return true
		}
	}
	return false
}

func (ball *Ball) Clearcheckpoint() {
	ball.Coord = nil
	ball.Checkpoints.Init()
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

func (balls *All_ball) Get_ballbyid_tomagnet(tab [3]int64) *list.List {
	list_tmp := list.New()
	var eb1 *list.Element
	var eb2 *list.Element
	var eb3 *list.Element

	for i := 0; i < 3; i++ {
		eball := balls.Blist.Front()
		for eball != nil && eball.Value.(*Ball).Id_ball != tab[i] {
			eball = eball.Next()
		}
		for eball.Value.(*Ball).Possessed != nil || (eball == eb1 || eball == eb2 || eball == eb3) {
			eball = eball.Next()
			if eball == nil {
				eball = balls.Blist.Front()
			}
		}
		switch i {
		case 0:
			eb1 = eball
		case 1:
			eb2 = eball
		case 2:
			eb3 = eball
		}
	}
	list_tmp.PushBack(eb1)
	list_tmp.PushBack(eb2)
	list_tmp.PushBack(eb3)
	return list_tmp
}

/* Apply the function Get_checkpointlist all ballons */
func (Lst_ball *All_ball) Create_checkpoint(Lst_wd *owm.All_data) error {
	var station owm.Weather_data

	station = Lst_wd.Get_Paris()
	Lst_ball.Lock()
	defer Lst_ball.Unlock()
	fmt.Println(station)
	eball := Lst_ball.Blist.Front()
	for eball != nil {
		eball.Value.(*Ball).Get_checkpointList(station)
		eball = eball.Next()
	}
	return nil
}

/* Give a next checkpoint ball and removes the previous */
func (Lst_ball *All_ball) Move_ball() (err error) {
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
					err = errors.New("next coord not found")
					return err
				}
			}
			ball.Stats.NbrKm += ball.AddStatsDistance(statsCoord.Lon, statsCoord.Lat)
		}
		elem.Value = ball
		elem = elem.Next()
	}
	return nil
}

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

func (Ball *Ball) GetItinerary() (int32, *list.List) {
	return 0, nil
}

func getIdMessageMax(idBall int64, base *db.Env) int32 {
	var IdMax int32
	err := base.Transact(base.Db, func(tx *sql.Tx) error {
		var err error
		stm, err := tx.Prepare("select id from message where id = (select max(id) from message) and containerid = $1;")
		checkErr(err)
		rs, err := stm.Query(idBall)
		checkErr(err)
		for rs.Next() {
			rs.Scan(&IdMax)
		}
		return err
	})
	checkErr(err)
	return IdMax
}

/*
 FUNCTION insertContainer(
	$1 idcreatorc integer,
	$2 latitudec integer,
	$3 longitudec integer,
	$4 device integer,
	$5 directionc float,
	$6 speedc float,
	$7 title text,
	$8 idx integer)
*/
func (Lst_ball *All_ball) InsertBallon(newBall *Ball, base *db.Env) (bool, error) {
	fmt.Printf("Insert  Id User %v | IdBall %v | \n", newBall.Creator.Value.(*users.User).Id, newBall.Id_ball)
	var err error
	var executed bool
	err = base.Transact(base.Db, func(tx *sql.Tx) error {
		stm, err := tx.Prepare("SELECT insertContainer($1, $2, $3, $4, $5, $6, $7 , $8)")
		checkErr(err)
		rs, err := stm.Query(newBall.Creator.Value.(*users.User).Id,
			newBall.Coord.Value.(Coordinate).Lat,
			newBall.Coord.Value.(Coordinate).Lon,
			3, //newBall.Creator.Value.(*users.User).Device.Front().Value.(users.Device).IdMobile,
			newBall.Wind.Degress,
			newBall.Wind.Speed,
			newBall.Title,
			newBall.Id_ball)
		for rs.Next() {
			var idC int
			err = rs.Scan(&idC)
			checkErr(err)
			err = Lst_ball.InsertMessages(newBall.Messages, idC, base)
			checkErr(err)
		}
		checkErr(err)
		return err
	})
	executed = true
	return executed, err
}

func (Lb *All_ball) Update_balls(ABalls *All_ball, base *db.Env) {
	i := 0
	fmt.Println("\x1b[31;1m coucou update\x1b[0m")
	fmt.Printf("%v Id Max", ABalls.Id_max)
	for e := ABalls.Blist.Front(); e != nil; e = e.Next() {

		if e.Value.(*Ball).Edited == true && e.Value.(*Ball).Id_ball <= ABalls.Id_max {

			idBall := e.Value.(*Ball).Id_ball
			idMessageMax := getIdMessageMax(idBall, base)
			j := 0
			for f := e.Value.(*Ball).Messages.Front(); f != nil; f = f.Next() {
				fmt.Printf(" %v id ballon update", f.Value.(Message).Id)
				if f.Value.(Message).Id > idMessageMax {
					err := base.Transact(base.Db, func(tx *sql.Tx) error {
						stm, err := tx.Prepare("INSERT INTO message(content, containerid, device_id) VALUES ($1, $2, $3)")
						checkErr(err)
						_, err = stm.Query(f.Value.(Message).Content, idBall, 2)
						j++
						checkErr(err)
						return err
					})
					checkErr(err)
				}
			}
		} else {

			Lb.InsertBallon(e.Value.(*Ball), base)
		}
		i++
	}
}

func (Lst_ball *All_ball) InsertMessages(messages *list.List, idBall int, base *db.Env) (err error) {
	i := 0
	for e := messages.Front(); e != nil; e = e.Next() {
		err = base.Transact(base.Db, func(tx *sql.Tx) error {
			stm, err := tx.Prepare("INSERT INTO message(content, containerid, device_id) VALUES ($1, $2, $3)")
			checkErr(err)
			_, err = stm.Query(e.Value.(Message).Content, idBall, 2)
			i++
			checkErr(err)
			return err
		})
	}
	return err
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
		panic(err)
	}
}

func (Lb *All_ball) GetFollowers(idBall int, Db *sql.DB, Ulist *list.List) *list.List {
	lstFollow := list.New()
	var err error
	rows, err := Db.Query("SELECT id_user FROM \"user\" AS userWibo LEFT OUTER JOIN followed ON (followed.iduser = userWibo.id_user)  WHERE followed.container_id = $1;", idBall)
	checkErr(err)
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
* RETURNS TABLE(idballon integer, titlename varchar(255), idtype integer, direction numeric, speedcont integer, creationdate date, deviceid integer, locationcont text)
 */

func (Lb *All_ball) GetListBallsByUser(userE *list.Element, base *db.Env, Ulist *list.List) *list.List {
	lBallon := list.New()
	var err error
	err = base.Transact(base.Db, func(tx *sql.Tx) error {
		var errT error
		stm, errT := tx.Prepare("SELECT getContainersByUserId($1)")
		checkErr(errT)
		rows, err := stm.Query(userE.Value.(*users.User).Id)
		checkErr(errT)
		for rows.Next() {
			var infoCont string
			err = rows.Scan(&infoCont)
			checkErr(err)
			result := strings.Split(infoCont, ",")
			idBall := GetIdBall(result[0])
			tempCord := GetCord(result[7])
			possessed := GetWhomGotBall(idBall, Ulist, base.Db)
			lBallon.PushBack(
				&Ball{
					Title:       result[1],
					Date:        GetDateFormat(result[5]),
					Checkpoints: tempCord,
					Coord:       tempCord.Front(),
					Wind:        GetWin(result[3], result[4]),
					Messages:    Lb.GetMessagesBall(idBall, base.Db),
					Followers:   Lb.GetFollowers(idBall, base.Db, Ulist),
					Possessed:   <-possessed,
					Creator:     userE})
			checkErr(err)
		}
		return err
	})
	checkErr(err)
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
	lMasterBall := new(All_ball)
	lMasterBall.Blist = list.New()
	i := 0
	for e := LstU.Ulist.Front(); e != nil; e = e.Next() {
		lMasterBall.Blist = Lb.GetListBallsByUser(e, base, LstU.Ulist)
		i++
	}
	Lb.InsertListBallsFollow(lMasterBall.Blist, LstU.Ulist, base)
	Lb.Blist = lMasterBall.Blist
	return nil
}
