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
	"github.com/Wibo/src/users"
	_ "github.com/lib/pq"
	"math"
	//	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
** Structure message qui contient le contenu de la string sa taille et son type
** Pour le moment seul le type 1 est pris en compte pour le texte.
** Le premier message de la liste est le message de creation du ballon.
 */
type Lst_msg struct {
	Content string
	Size    int
	Type_   int
}

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

/*
** Structure Checkpoints classique avec les coordonnees du checkpoint
** et la date du checkpoints.
** ! Il n'y a pas encore de gestion d'historique des checkpoints !
 */
type Checkpoints struct {
	Coord Coordinates
	Date  time.Time
}

type Wind struct {
	Speed   float64
	Degress float64
}

/*
** Interface Ball, Name contient le titre du ballon.
** Coord est un element de la liste Checkpoints.
** ! Wind concerne les vents qui sont appliques sur le ballon !
** ! (Non utilise pour seule station consulte, Paris) !
** Lst_msg est une liste de message de type list.Element.Value.(Lst_msg)
** Date est la date de creation du ballon.
** Checkpoints est la liste des checkpoints sur les prochaines 3 heures a
** intervale de 5 minutes.
** Possessed est l'user qui possede le ballon actuellement.
** List_follow est une liste d'utilisateur suivant le ballon de type
** list.Element.Value.(users.User). Un user sera un lien vers le meme user
** provenant de la liste d'utilisateur principale afin d'eviter les datas
** multiples.
** Creator est l'user qui a creer le ballon.
 */
type Ball struct {
	//idBall      int64
	Name        string
	Position    Coordinates
	Wind        Wind
	Lst_msg     *list.List
	Date        time.Time
	Checkpoints *list.List
	Coord       *list.Element
	Possessed   *users.User
	List_follow *list.List
	Creator     *users.User
}

/*
** All_ball est la structure principale des ballons.
** Elle contient un Mutex et une liste de type list.Element.Value.(Ball)
 */
type All_ball struct {
	sync.RWMutex
	Lst *list.List
}

/* Print_list_checkpoints print la liste de checkpoints d'un ballon */
func (ball Ball) Print_list_checkpoints() {
	elem := ball.Checkpoints.Front()

	for elem != nil {
		fmt.Println(elem.Value.(Checkpoints))
		elem = elem.Next()
	}
}

/*
** Get_checkpointList creer la liste de checkpoints d'un ballon.
** Dans le cadre de la beta, il verifie les coordonnees du ballon pour le
** forcer a rester dans Paris.
 */
func (elem Ball) Get_checkpointList(station owm.Weather_data) (test Ball) {
	r_world := 6371000.0
	var tmp_coord Coordinates
	var calc_coord Coordinates
	var checkpoint Coordinates

	speed := station.Wind.Speed * 300.00
	dir := station.Wind.Degress*10 + 180
	if dir >= 360 {
		dir -= 360
	}
	dir = dir * (math.Pi / 180.0)
	checkpoint.Longitude = elem.Coord.Value.(Checkpoints).Coord.Longitude
	checkpoint.Latitude = elem.Coord.Value.(Checkpoints).Coord.Latitude
	elem.Checkpoints = elem.Checkpoints.Init()
	for i := 0; i < 35; i++ {
		tmp_coord.Longitude = checkpoint.Longitude * (math.Pi / 180.0)
		tmp_coord.Latitude = checkpoint.Latitude * (math.Pi / 180.0)
		calc_coord.Latitude = math.Asin(math.Sin(tmp_coord.Latitude)*math.Cos(speed/r_world) + math.Cos(tmp_coord.Latitude)*math.Sin(speed/r_world)*math.Cos(dir))
		calc_coord.Longitude = tmp_coord.Longitude + math.Atan2(math.Sin(dir)*math.Sin(speed/r_world)*math.Cos(tmp_coord.Latitude), math.Cos(speed/r_world)-math.Sin(tmp_coord.Latitude)*math.Sin(calc_coord.Latitude))
		calc_coord.Latitude = 180 * calc_coord.Latitude / math.Pi
		calc_coord.Longitude = 180 * calc_coord.Longitude / math.Pi
		if calc_coord.Latitude < 2.10 {
			checkpoint.Latitude = 2.60
		} else if calc_coord.Latitude > 2.60 {
			checkpoint.Latitude = 2.10
		} else {
			checkpoint.Latitude = calc_coord.Latitude
		}
		if calc_coord.Longitude < 48.72 {
			checkpoint.Longitude = 49.02
		} else if calc_coord.Longitude > 49.02 {
			checkpoint.Longitude = 48.72
		} else {
			checkpoint.Longitude = calc_coord.Longitude
		}
		elem.Checkpoints.PushBack(Checkpoints{checkpoint, time.Now()})
	}
	/* CECI EST UN TEST DE FONCTIONNALITE */
	elem.Print_list_checkpoints()
	/* FIN DU TEST */
	return elem
}

/*
** Create checkpoint applique a tous les ballon, la nouvelle liste de
** checkpoints qui leur correspondent. Cette fonction est appele toutes
** les 3 heures, quand la liste de checkpoints d'un ballon est vide.
** !Apres la beta, faire une fonction qui va regarder les 3 stations les
** plus proches de la position actuelle du ballon pour en definir un vecteur
** de vent. Ce membre de fonction sera ainsi appele dans le for.
 */
func (Lst_ball *All_ball) Create_checkpoint(Lst_wd *owm.All_data) error {
	var station owm.Weather_data
	for _, elem := range Lst_wd.Tab_wd {
		if elem.Station_name == "Paris" {
			station = elem
			break
		}
	}
	Lst_ball.Lock()
	defer Lst_ball.Unlock()
	fmt.Println(station)
	elem := Lst_ball.Lst.Front()
	for elem != nil {
		elem.Value = elem.Value.(Ball).Get_checkpointList(station)
		elem = elem.Next()
	}
	return nil
}

/*
** Move_ball est appelle toutes les 5 minutes pour changer les coordonnees de
** tous les ballons.
 */
func (Lst_ball *All_ball) Move_ball() (err error) {
	Lst_ball.Lock()
	defer Lst_ball.Unlock()
	elem := Lst_ball.Lst.Front()

	for elem != nil {
		ball := elem.Value.(Ball)
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

/* Add_new_ballon va ajouter un ballon suite a une requete client. */
func (Lst_ball *All_ball) Add_new_ballon(new_ball Ball) {
	Lst_ball.Lst.PushBack(new_ball)
	return
}

func (Lst_ball *All_ball) Update_new_ballon(upd_ball *Ball) {
	return
}

func (Lst_ball *All_ball) Print_all_balls() {
	// Print list to debug
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
* GetListBallsByUser
* getContainersByUserId is a native psql function with
* RETURNS TABLE(idballon integer, titlename varchar(255), idtype integer, direction numeric, speedcont integer, creationdate date, deviceid integer, locationcont text)
 */
func (Lb *All_ball) GetListBallsByUser(idUser int64, Db *sql.DB) *list.List {

	var err error
	lBallon := list.New()
	stm, err := Db.Prepare("SELECT getContainersByUserId($1)")
	checkErr(err)
	rows, err := stm.Query(idUser)
	checkErr(err)
	// regex to find words
	//r, err := regexp.Compile(`[:print:]\w+`)// getName
	for rows.Next() {
		var infoCont string
		err = rows.Scan(&infoCont)
		checkErr(err)
		result := strings.Split(infoCont, ",")
		for i := range result {
			fmt.Printf("%v id ballon\n", GetIdBall(result[0]))
			fmt.Printf("%v and position %v\n", result[i], i)
			lBallon.PushBack(Ball{Name: result[1], Date: GetDateFormat(result[5]), Position: GetCord(result[7]),
				Wind: GetWin(result[5], result[4])})
		}
	}
	return lBallon
}

func GetDateFormat(qdate string) (fdate time.Time) {
	// TODO Choose a date format layout
	fdate, err := time.Parse("2006-01-02", qdate)
	checkErr(err)
	return fdate
}

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

func GetCord(position string) Coordinates {

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
	long, _ := strconv.ParseFloat(point[0], 4)
	lat, _ := strconv.ParseFloat(point[1], 4)
	return (Coordinates{Longitude: long, Latitude: lat})
}

func GetWin(speed string, direction string) Wind {
	sf, _ := strconv.ParseFloat(speed, 4)
	df, _ := strconv.ParseFloat(direction, 4)
	return (Wind{Speed: sf, Degress: df})
}

func GetMessagesBall(idBall int) *list.List {
	return nil
}

/**
* get all ball from database and associeted
* the creator, possessord and followers.
**/
func (Lb *All_ball) Get_balls(LstU *users.All_users, Db *sql.DB) error {
	lMasterBall := list.New()
	i := 0
	for e := LstU.Lst_users.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v | %v | %v \n", e.Value.(users.User).Id_user, e.Value.(users.User).Login)
		i++
	}
	Lb.Lst = Lb.Lst.Init()
	Lb.Lst.PushFrontList(lMasterBall)
	fmt.Println(lMasterBall)
	return nil
}
