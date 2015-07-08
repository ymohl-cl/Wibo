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
	Longitude float64
	Latitude  float64
}

type Wind struct {
	Speed   float64
	Degress float64
}

/* Date is a timestamp's creation to Ball */
type Ball struct {
	name        string
	coord       Coordinates
	Wind        Wind
	Lst_msg     []Lst_msg
	date        int
	Checkpoints []Checkpoints
	Possessed   *users.User
	List_follow []users.All_users
	Creator     *users.User
}

type All_ball struct {
	Lst_ball []Ball
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

// Add the user in the ball on creation ball list
func (Lst_ball *All_ball) Get_balls(Lst_users *users.All_users) error {
	// Get all information
	return nil
}
