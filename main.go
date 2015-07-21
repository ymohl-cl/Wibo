//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  main.go                                            :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  By: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  Created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  Updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

// ** Creer un nouveau fil pour effectuer un membre de fonction qui va creer un fichier
// ** avec le nom de la date et heure de creation et la sauvegarde as log.

package main

import (
	"ballon"
	"container/list"
	"fmt"
	"net/http"
	"owm"
	"request"
	"sock"
	"time"
	"users"
)

func Manage_goroutines(Tab_wd *owm.All_data, Lst_ball *ballon.All_ball) {
	channelfuncweatherdata := make(chan bool)
	channelfuncmoveball := make(chan bool)
	defer close(channelfuncmoveball)
	defer close(channelfuncweatherdata)

	go func() {
		for {
			select {
			case <-time.After(3 * time.Hour):
				channelfuncweatherdata <- true
			case <-time.After(5 * time.Minute):
				channelfuncmoveball <- true
			}
		}
	}()

	for {
		select {
		case <-channelfuncweatherdata:
			{
				err := Tab_wd.Update_weather_data()
				if err != nil {
					fmt.Println(err)
				} else {
					Tab_wd.Print_weatherdata()
				}
				err = Lst_ball.Create_checkpoint(Tab_wd)
				if err != nil {
					fmt.Println(err)
				} else {
					Lst_ball.Print_all_balls()
				}
			}
		case <-channelfuncmoveball:
			{
				fmt.Println("move coord on next checkpoint")
				err := Lst_ball.Move_ball()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func init_all(Tab_wd *owm.All_data, Lst_users *users.All_users, Lst_ball *ballon.All_ball) error {
	// Get first array data
	err := Tab_wd.Update_weather_data()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Tab_wd.Print_weatherdata()
	}

	// Get first list user
	err = Lst_users.Get_users()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Lst_users.Print_users()
	}

	// Get first list ballon with their follower

	Lst_ball.Lst = list.New()
	err = Lst_ball.Get_balls(Lst_users)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Lst_ball.Print_all_balls()
	}
	// to test with one ball
	tmp_lst := list.New()
	var check_test ballon.Checkpoints
	check_test.Coord.Longitude = 48.833086
	check_test.Coord.Latitude = 2.316055

	check_test.Date = 0
	var my_ball = ballon.Ball{"toto", nil, ballon.Wind{}, list.New(), 0, list.New(), nil, users.All_users{}, nil, nil}
	my_ball.Coord = tmp_lst.PushBack(check_test)
	fmt.Println("Debug to get checkpoint coord")
	fmt.Println(my_ball.Coord.Value)

	Lst_ball.Add_new_ballon(my_ball)

	// Get first list checkpoints ball
	err = Lst_ball.Create_checkpoint(Tab_wd)
	if err != nil {
		fmt.Println(err)
	} else {
		Lst_ball.Print_all_balls()
	}

	return nil
}

/*
** Les requetes sont utilise que pour recuperer la positon
** des ballons autour de la position recu.
** Si il y a une modification sur un ballon, envoyer une
** requetes a toutes les client encore ON en HTTP ou en socket
** si elle est encore ouverte.
** Les socket sont utilisees pour tous les autres types
** de communications.
 */

func main() {
	Tab_wd := new(owm.All_data)
	Lst_users := new(users.All_users)
	Lst_ball := new(ballon.All_ball)

	err := init_all(Tab_wd, Lst_users, Lst_ball)
	if err != nil {
		return
	}
	go Manage_goroutines(Tab_wd, Lst_ball)

	request.Init_handle_request()
	go http.ListenAndServe(":8080", nil)
	go sock.Listen(Lst_users)

	for {
		fmt.Println("manage server")
		time.Sleep(time.Second * 60)
	}
}
