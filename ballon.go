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

/* Type is define to know a style/type message, test:1 / photo:2 / other:3
/* and by default:0 */
type Lst_msg {
	content  string
	size     int
	type_    int
}

type Checkpoints struct {
	Coord     Coordinates
	date      int
}

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

type Wind struct {
	Speed   float64
	Degress float64
}

/* First message is message creation and date is a timestamp creation,
/* Checkpoints is a list checkpoints in order to most old at less old */
type All_ball struct {
	name         string
	coord        Coordinates
	Wind         Wind
	Lst_msg      Lst_msg
	date         int
	Checkpoints  Checkpoints
}

func (Lst_ball *All_ball) Get_ball(Lst_ball *All_ball) error {
	Lst_ball := new (All_ball)
	// Get all information

	return Lst_ball
}
