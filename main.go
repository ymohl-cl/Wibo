//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  main.go                                            :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  by: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

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

/*
** Manage_goroutines va gerer les differentes processus endormis
** qui effectueront des taches tous les X temps.
** Synchronisation des datas des vents toutes les 3 heures
** Deplacement des ballons grace au pre calcul de checkpoint toutes les 5 minutes
** Creation de checkpoint toutes les 3 heures.
** !! Viendra la synchronisation du cache dans la base de donnee tous les X temps !!
 */
func Manage_goroutines(Tab_wd *owm.All_data, Lst_ball *ballon.All_ball) {
	channelfuncweatherdata := make(chan bool)
	channelfuncmoveball := make(chan bool)
	defer close(channelfuncmoveball)
	defer close(channelfuncweatherdata)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			channelfuncmoveball <- true
		}
	}()
	go func() {
		for {
			time.Sleep(3 * time.Hour)
			channelfuncweatherdata <- true
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
				err := Lst_ball.Move_ball()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

/*
** Init_all initialise toutes les datas en recuperant celles presentes dans la base de donnee
** 1: On recupere les datas des vents.
** 2: On recupere la liste des utilisateurs dans la base de donnee.
** 3: On recupere la liste des ballons dans la base de donnee et on y attache les users concernes par le ballon
** 4: On cree la liste des checkpoints pour chaque ballon.
 */
func Init_all(Tab_wd *owm.All_data, Lst_users *users.All_users, Lst_ball *ballon.All_ball) error {
	err := Tab_wd.Update_weather_data()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Tab_wd.Print_weatherdata()
	}
	Lst_users.Lst_users = list.New()
	err = Lst_users.Get_users()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Lst_users.Print_users()
	}
	Lst_ball.Lst = list.New()
	err = Lst_ball.Get_balls(Lst_users)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		Lst_ball.Print_all_balls()
	}

	/* CREER UN BALLON POUR FAIRE DES TESTS */
	tmp_lst := list.New()
	var check_test ballon.Checkpoints
	check_test.Coord.Longitude = 48.833086
	check_test.Coord.Latitude = 2.316055
	check_test.Date = time.Now()
	var my_ball = ballon.Ball{Id_ball: 42, Name: "toto", Coord: nil, Wind: ballon.Wind{}, Lst_msg: list.New(), Date: time.Now(), Checkpoints: list.New(), Possessed: nil, List_follow: list.New(), Creator: nil}
	my_ball.Coord = tmp_lst.PushBack(check_test)
	fmt.Println("Debug to get checkpoint coord")
	fmt.Println(my_ball.Coord.Value)
	Lst_ball.Add_new_ballon(my_ball)
	/* FIN DE LA CREATION DEBALLON POUR TEST */

	err = Lst_ball.Create_checkpoint(Tab_wd)
	if err != nil {
		fmt.Println(err)
	} else {
		Lst_ball.Print_all_balls()
	}
	return nil
}

/*
** Les 3 datats essentielle sont instancie dans le main
** Tab_wd contient toutes les WEATHER DATA
** Lst_users contients tous les utilisateurs
** La liste ballon contient tous les ballons
** Ont initialise nos 3 datas
** On lance les goroutines
** On initialise la reception de requete http (Pour le moment non utilise)
** On initialise la reception de requete socket (Actuellement utilise)
** Et on a une boucle for qui empeche la fermeture du programme.
 */
func main() {
	Tab_wd := new(owm.All_data)
	Lst_users := new(users.All_users)
	Lst_ball := new(ballon.All_ball)

	err := Init_all(Tab_wd, Lst_users, Lst_ball)
	if err != nil {
		return
	}
	go Manage_goroutines(Tab_wd, Lst_ball)

	request.Init_handle_request()
	go http.ListenAndServe(":8080", nil)
	go sock.Listen(Lst_users, Lst_ball, Tab_wd)

	for {
		time.Sleep(time.Second * 60)
	}
}
