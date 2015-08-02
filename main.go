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
	var check_test0 ballon.Checkpoints
	check_test0.Coord.Longitude = 48.833086
	check_test0.Coord.Latitude = 2.316055
	check_test0.Date = time.Now()
	var check_test1 ballon.Checkpoints
	check_test1.Coord.Longitude = 48.833586
	check_test1.Coord.Latitude = 2.316065
	check_test1.Date = time.Now()
	var check_test2 ballon.Checkpoints
	check_test2.Coord.Longitude = 48.833368
	check_test2.Coord.Latitude = 2.316059
	check_test2.Date = time.Now()
	var check_test3 ballon.Checkpoints
	check_test3.Coord.Longitude = 48.833286
	check_test3.Coord.Latitude = 2.316903
	check_test3.Date = time.Now()
	var check_test4 ballon.Checkpoints
	check_test4.Coord.Longitude = 48.833986
	check_test4.Coord.Latitude = 2.316045
	check_test4.Date = time.Now()
	mmp2 := list.New()
	mmp := list.New()
	var message0 ballon.Lst_msg
	message0.Id_Message = 0
	message0.Size = 68
	message0.Content = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean com"
	message0.Type_ = 1
	var message1 ballon.Lst_msg
	message1.Id_Message = 0
	message1.Size = 68
	message1.Content = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean com"
	message1.Type_ = 1
	var message2 ballon.Lst_msg
	message2.Id_Message = 2
	message2.Size = 68
	message2.Content = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean com"
	message2.Type_ = 1

	mmp.PushBack(message0)
	mmp2.PushBack(message0)
	mmp2.PushBack(message1)
	mmp2.PushBack(message2)

	ball0 := new(ballon.Ball)
	ball0.Id_ball = 0
	ball0.Name = "toto"
	ball0.Coord = tmp_lst.PushBack(check_test0)
	ball0.Wind = ballon.Wind{}
	ball0.Lst_msg = mmp
	ball0.Date = time.Now()
	ball0.Checkpoints = list.New()
	ball0.Possessed = nil
	ball0.List_follow = list.New()
	ball0.Creator = nil
	Lst_ball.Lst.PushBack(ball0)

	ball1 := new(ballon.Ball)
	ball1.Id_ball = 1
	ball1.Name = "tata"
	ball1.Coord = tmp_lst.PushBack(check_test1)
	ball1.Wind = ballon.Wind{}
	ball1.Lst_msg = mmp
	ball1.Date = time.Now()
	ball1.Checkpoints = list.New()
	ball1.Possessed = nil
	ball1.List_follow = list.New()
	ball1.Creator = nil
	Lst_ball.Lst.PushBack(ball1)

	ball2 := new(ballon.Ball)
	ball2.Id_ball = 2
	ball2.Name = "tutu"
	ball2.Coord = tmp_lst.PushBack(check_test2)
	ball2.Wind = ballon.Wind{}
	ball2.Lst_msg = mmp
	ball2.Date = time.Now()
	ball2.Checkpoints = list.New()
	ball2.Possessed = nil
	ball2.List_follow = list.New()
	ball2.Creator = nil
	Lst_ball.Lst.PushBack(ball2)

	ball3 := new(ballon.Ball)
	ball3.Id_ball = 3
	ball3.Name = "tete"
	ball3.Coord = tmp_lst.PushBack(check_test3)
	ball3.Wind = ballon.Wind{}
	ball3.Lst_msg = mmp
	ball3.Date = time.Now()
	ball3.Checkpoints = list.New()
	ball3.Possessed = nil
	ball3.List_follow = list.New()
	ball3.Creator = nil
	Lst_ball.Lst.PushBack(ball3)

	ball4 := new(ballon.Ball)
	ball4.Id_ball = 4
	ball4.Name = "tyty"
	ball4.Coord = tmp_lst.PushBack(check_test4)
	ball4.Wind = ballon.Wind{}
	ball4.Lst_msg = mmp
	ball4.Date = time.Now()
	ball4.Checkpoints = list.New()
	ball4.Possessed = nil
	ball4.List_follow = list.New()
	ball4.Creator = nil
	Lst_ball.Lst.PushBack(ball4)
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
