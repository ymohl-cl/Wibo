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
	"errors"
	"fmt"
	"math"
	"owm"
	"sync"
	"time"
	"users"
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
	Name        string
	Coord       *list.Element
	Wind        Wind
	Lst_msg     *list.List
	Date        time.Time
	Checkpoints *list.List
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

/* Update_new_ballon va mettre a jour un ballon suite a une requete client. */
func (Lst_ball *All_ball) Update_new_ballon(upd_ball *Ball) {
	return
}

/* Print_all_ball print la liste de tous les ballons, utile pour debeuguer. */
func (Lst_ball *All_ball) Print_all_balls() {
	return
}

/*
** Get_balls recupere tous les ballons de la base de donnee.
** Elle rattache egalement a chaque ballon tous les utilisateurs associes.
 */
func (Lst_ball *All_ball) Get_balls(Lst_users *users.All_users) error {
	Lst_ball.Lock()
	//  ...
	//  traitement
	//  ...
	Lst_ball.Unlock()
	return nil
}
