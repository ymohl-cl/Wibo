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
	"fmt"
	"owm"
	"time"
)

//func Manage_goroutines(Lst_wd *owm.All_data) {
//	channelfuncweatherdata := make(chan bool)
//	channelfunccheckpointball := make(chan bool)
//	defer close(channelfunccheckpointball)
//	defer close(channelfuncweatherdata)
//
//	go func() {
//		for {
//			select {
//			case <-time.After(5 * time.Second):
//				channelfuncweatherdata <- true
//			}
//		}
//	}()
//	for {
//		select {
//		case <-channelfuncweatherdata:
//			{
//				err := Lst_wd.Update_weather_data()
//				if err != nil {
//					fmt.Println(err)
//				} else {
//					channelfunccheckpointball <- true
//				}
//			}
//		case <-channelfunccheckpointball: /*go create_checkpointball() */
//			Lst_wd.Print_weatherdata()
//		}
//	}
//	fmt.Println("End manage_goroutines()\n")
//}

func Manage_goroutines(Lst_wd *owm.All_data) {
	channelfuncweatherdata := make(chan bool)
	channelfunccheckpointball := make(chan bool)
	defer close(channelfunccheckpointball)
	defer close(channelfuncweatherdata)

	go func() {
		for {
			select {
			case <-time.After(30 * time.Second):
				channelfuncweatherdata <- true
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
		}
	}
	fmt.Println("End manage_goroutines()\n")
}

func main() {
	Lst_wd := new(owm.All_data)
	go Manage_goroutines(Lst_wd)

	for {
		fmt.Println("manage server")
		time.Sleep(time.Minute * 1)
	}
}
