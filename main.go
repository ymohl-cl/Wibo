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
	"fmt"
	"net/http"
	"owm"
	"request"
	"time"
	"users"
)

func Manage_goroutines(Lst_wd *owm.All_data, Lst_ball *ballon.All_ball) {
	channelfuncweatherdata := make(chan bool)
	channelfunccheckpointball := make(chan bool)
	defer close(channelfunccheckpointball)
	defer close(channelfuncweatherdata)

	go func() {
		for {
			select {
			case <-time.After(3600 * time.Second):
				channelfuncweatherdata <- true
				channelfunccheckpointball <- true
			}
		}
	}()

	for {
		select {
		case <-channelfuncweatherdata:
			{
				err := Lst_wd.Update_weather_data()
				if err != nil {
					fmt.Println(err)
				} else {
					Lst_wd.Print_weatherdata()
				}
			}
		case <-channelfunccheckpointball:
			{
				err := Lst_ball.Create_checkpoint(Lst_wd)
				if err != nil {
					fmt.Println(err)
				} else {
					Lst_ball.Print_all_balls()
				}
			}
		}
	}
	fmt.Println("End manage_goroutines()\n")
}

func init_all(Lst_wd *owm.All_data, Lst_users *users.All_users, Lst_ball *ballon.All_ball) error {
	err := Lst_wd.Update_weather_data()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Lst_wd.Print_weatherdata()
	}

	err = Lst_users.Get_users()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Lst_users.Print_users()
	}

	err = Lst_ball.Get_balls(Lst_users)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Lst_ball.Print_all_balls()
	}
	return nil
}

func main() {
	Lst_wd := new(owm.All_data)
	Lst_users := new(users.All_users)
	Lst_ball := new(ballon.All_ball)

	err := init_all(Lst_wd, Lst_users, Lst_ball)
	if err != nil {
		return
	}
	go Manage_goroutines(Lst_wd, Lst_ball)
	request.Init_handle_request()

	go http.ListenAndServe(":8080", nil)

	for {
		fmt.Println("manage server")
		time.Sleep(time.Second * 60)
	}
}
