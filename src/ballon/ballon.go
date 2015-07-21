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
	"users"
)

/* Type is define to know a style/type message, test:1 / photo:2 / other:3
/* First message is message creation */
/* and by default:0 */
type Lst_msg struct {
	Content string
	Size    int
	Type_   int
}

/* Checkpoints is a list checkpoints in order to most old at less old */
type Checkpoints struct {
	Coord Coordinates
	Date  int
}

type Coordinates struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type Wind struct {
	Speed   float64
	Degress float64
}

/* Date is a timestamp's creation to Ball */
// Add to ball, history checkpoints with first and last elem of list checkpoints
// every 3 hours. After, checkpoints list, is refresh too every 3 hours.
type Ball struct {
	Name        string
	Coord       *list.Element
	Wind        Wind
	Lst_msg     *list.List
	Date        int
	Checkpoints *list.List
	Possessed   *users.User
	List_follow users.All_users
	Creator     *users.User
	Next        *Ball
}

type All_ball struct {
	Lst *list.List
}

// Ajouter a un ballon un attribut history checkpoints qui
// prendra le premier et le dernier element des checkpoints
// list pour pouvoir constituer un itineraire leger.
func (elem Ball) Get_checkpointList(station owm.Weather_data) {
	r_world := 6371000.0
	var tmp_coord Coordinates
	var calc_coord Coordinates
	var checkpoint Coordinates

	fmt.Println("Infos ballon")
	fmt.Println("Coordonnees:")
	fmt.Println(elem.Coord.Value.(Checkpoints))
	speed := station.Wind.Speed * 300.00
	dir := station.Wind.Degress*10 + 180
	if dir >= 360 {
		dir -= 360
	}
	dir = dir * (math.Pi / 180.0)
	checkpoint.Longitude = elem.Coord.Value.(Checkpoints).Coord.Longitude
	checkpoint.Latitude = elem.Coord.Value.(Checkpoints).Coord.Latitude

	for i := 0; i < 36; i++ {

		tmp_coord.Longitude = checkpoint.Longitude * (math.Pi / 180.0)
		tmp_coord.Latitude = checkpoint.Latitude * (math.Pi / 180.0)

		calc_coord.Latitude = math.Asin(math.Sin(tmp_coord.Latitude)*math.Cos(speed/r_world) + math.Cos(tmp_coord.Latitude)*math.Sin(speed/r_world)*math.Cos(dir))

		calc_coord.Longitude = tmp_coord.Longitude + math.Atan2(math.Sin(dir)*math.Sin(speed/r_world)*math.Cos(tmp_coord.Latitude), math.Cos(speed/r_world)-math.Sin(tmp_coord.Latitude)*math.Sin(calc_coord.Latitude))

		calc_coord.Latitude = 180 * calc_coord.Latitude / math.Pi

		calc_coord.Longitude = 180 * calc_coord.Longitude / math.Pi

		fmt.Println("checkpoint new")
		fmt.Println("Wait ... latitude")
		if calc_coord.Latitude < 2.10 {
			fmt.Println(calc_coord.Latitude)
			fmt.Println("modiff lat 2.60")
			checkpoint.Latitude = 2.60
		} else if calc_coord.Latitude > 2.60 {
			fmt.Println(calc_coord.Latitude)
			fmt.Println("modiff lat 2.10")
			checkpoint.Latitude = 2.10
		} else {
			fmt.Println(calc_coord.Latitude)
			fmt.Println("not modiff lat")
			checkpoint.Latitude = calc_coord.Latitude
		}

		fmt.Println("Wait ... longitude")
		if calc_coord.Longitude < 48.72 {
			fmt.Println(calc_coord.Longitude)
			fmt.Println("modiff lon 49.02")
			checkpoint.Longitude = 49.02
		} else if calc_coord.Longitude > 49.02 {
			fmt.Println(calc_coord.Longitude)
			fmt.Println("modiff lon 48.72")
			checkpoint.Longitude = 48.72
		} else {
			fmt.Println(calc_coord.Longitude)
			checkpoint.Longitude = calc_coord.Longitude
			fmt.Println("not modiff lon")
		}

		fmt.Printf("checkpoint: %dmin\n: ", (i*5 + 5))
		fmt.Println(checkpoint)
		elem.Checkpoints.PushBack(checkpoint)
	}
	return
}

func (Lst_ball *All_ball) Create_checkpoint(Lst_wd *owm.All_data) error {
	// Get station to calcul checkpoint. Later, this, will be a member function that elem.
	var station owm.Weather_data
	for _, elem := range Lst_wd.Tab_wd {
		if elem.Station_name == "Paris" {
			station = elem
			break
		}
	}
	// ------------------------------------------------------------------------------------
	fmt.Println(station)
	elem := Lst_ball.Lst.Front()
	fmt.Println("Get first:D")
	for elem != nil {
		fmt.Println("begin for")
		elem.Value.(Ball).Get_checkpointList(station)
		elem = elem.Next()
		fmt.Println("end for")
	}
	fmt.Println("Debug :D")
	return nil
}

func (Lst_ball *All_ball) Move_ball() (err error) {
	elem := Lst_ball.Lst.Front()
	for elem != nil {
		ball := elem.Value.(Ball)
		fmt.Println("Old coord")
		fmt.Println(ball.Coord.Value)
		new_coord := ball.Coord.Next()
		if new_coord != nil {
			fmt.Println("Next coord exit, new coord")
			fmt.Println(ball.Coord.Value)
			ball.Coord = new_coord
		} else {
			new_coord = ball.Checkpoints.Front()
			if new_coord != nil {
				fmt.Println("Next coord not exit, take liste checkpoints, new coord:")
				fmt.Println(ball.Coord.Value)
				ball.Coord = new_coord
			} else {
				err = errors.New("next coord not found")
				return err
			}
		}
	}
	return nil
}

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

func (Lst_ball *All_ball) Get_balls(Lst_users *users.All_users) error {
	//get all ball from database and associeted
	// the creator, possessord and followers.
	return nil
}
