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
	"github.com/Wibo/src/owm"
	"github.com/Wibo/src/usr"
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
type Ball struct {
	Name        string      `json:"name"`
	Coord       Coordinates `json:"coord"`
	Wind        Wind
	Lst_msg     list.List
	Date        int
	Checkpoints list.List
	Possessed   *users.User
	List_follow users.All_users
	Creator     *users.User
	Next        *Ball
}

type All_ball struct {
	Lst *list.List `json:"list"`
}

func (Lst_ball *All_ball) Create_checkpoint(Lst_wd *owm.All_data) error {
	// Create all checkpoint
	return nil
}

func (Lst_ball *All_ball) Add_new_ballon(new_ball *Ball) {
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
