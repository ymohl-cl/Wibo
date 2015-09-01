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
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Wibo/src/owm"
	"github.com/Wibo/src/protocol"
	"github.com/Wibo/src/users"
	//_ "github.com/lib/pq"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

/* Type is message type. Only type 1 is use now and described a text */
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

type Checkpoint struct {
	Coord Coordinate
	Date  time.Time
}

type Wind struct {
	Speed   float64
	Degress float64
}

type Ball struct {
	Id_ball     int64
	Title       string
	Coord       *list.Element
	Wind        Wind
	Messages    *list.List    /* Value: Message */
	Date        time.Time     /* creation date */
	Checkpoints *list.List    /* list checkpoints's ball a five minutes inteval */
	Possessed   *list.Element /* Value: (*users.User) */
	Followers   *list.List    /* Value: *list.Element.Value.(*users.User) */
	Creator     *list.Element /* Value: (*users.User) */
}

type All_ball struct {
	sync.RWMutex
	Blist  *list.List /* Value: *Ball */
	Id_max int64      /* Set by bdd and incremented by server */
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

func (ball *Ball) Check_nearbycoord(request *list.Element) bool {
	rlon := request.Value.(protocol.Request).Spec.(protocol.Position).Lon
	rlat := request.Value.(protocol.Request).Spec.(protocol.Position).Lat
	coord := ball.Coord.Value.(Checkpoint).Coord
	if coord.Lon < rlon+0.01 &&
		coord.Lon > rlon-0.01 &&
		coord.Lat < rlat+0.01 &&
		coord.Lat > rlat-0.01 {
		return true
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
		if calc_coord.Lat < 2.10 {
			checkpoint.Lat = 2.60
		} else if calc_coord.Lat > 2.60 {
			checkpoint.Lat = 2.10
		} else {
			checkpoint.Lat = calc_coord.Lat
		}
		if calc_coord.Lon < 48.72 {
			checkpoint.Lon = 49.02
		} else if calc_coord.Lon > 49.02 {
			checkpoint.Lon = 48.72
		} else {
			checkpoint.Lon = calc_coord.Lon
		}
		eball.Checkpoints.PushBack(Checkpoint{checkpoint, time.Now()})
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
		fmt.Println(user.Device)
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
* GetListBallsByUser
* getContainersByUserId is a native psql function with
* RETURNS TABLE(idballon integer, titlename varchar(255), idtype integer, direction numeric, speedcont integer, creationdate date, deviceid integer, locationcont text)
 */
func (Lb *All_ball) GetListBallsByUser(userl users.User, Db *sql.DB) *list.List {
	var err error
	lBallon := list.New()
	stm, err := Db.Prepare("SELECT getContainersByUserId($1)")
	checkErr(err)
	rows, err := stm.Query(userl.Id)
	checkErr(err)
	// regex to find words
	//r, err := regexp.Compile(`[:print:]\w+`)// getName
	for rows.Next() {
		var infoCont string
		err = rows.Scan(&infoCont)
		checkErr(err)
		result := strings.Split(infoCont, ",")
		lBallon.PushBack(
			Ball{Title: result[1],
				Date:  GetDateFormat(result[5]),
				Coord: nil, //Coord cannot use GetCord(result[7]), he muste
				// take *list.Element type
				Wind:     GetWin(result[3], result[4]),
				Messages: GetMessagesBall(GetIdBall(result[0]), Db), Creator: nil})
	}
	// Creator cannot == &userl. you must take Ulist and found User fo get *list.Element.
	return lBallon
}

/**
* GetDateFormat
* Set the string "Date psql format"  "YY-MM-DD HH:MM:SS.mls+00" format time.time
* return time.Time
 */

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
func GetCord(position string) Coordinate {
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
	return (Coordinate{Lon: long, Lat: lat})
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
func GetMessagesBall(idBall int, Db *sql.DB) *list.List {
	Mlist := list.New()
	stm, err := Db.Prepare("SELECT id AS containerId, content, id_type_m  FROM message WHERE containerid=($1) ORDER BY creationdate DESC")
	checkErr(err)
	rows, err := stm.Query(idBall)
	checkErr(err)
	for rows.Next() {
		var idm int
		var message string
		var idType int
		err = rows.Scan(&idm, &message, &idType)
		checkErr(err)
		Mlist.PushBack(Message{Content: message, Type: (int32)(idType), Id: int32(idm)})
	}
	return Mlist
}

/**
* get all ball from database and associeted
* the creator, possessord and followers.
**/
func (Lb *All_ball) Get_balls(LstU *users.All_users, Db *sql.DB) error {
	lMasterBall := list.New()
	i := 0
	for e := LstU.Ulist.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v | %v \n", e.Value.(users.User).Id, e.Value.(users.User).Login)
		lMasterBall.PushBackList(Lb.GetListBallsByUser(e.Value.(users.User), Db))
		i++
	}
	Lb.Blist.Init()
	Lb.Blist.PushBackList(lMasterBall)
	return nil
}
