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
	"Wibo/ballon"
	"Wibo/ballonwork"
	"Wibo/crontask"
	"Wibo/db"
	"Wibo/devices"
	"Wibo/owm"
	"Wibo/server"
	"Wibo/sock"
	"Wibo/users"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

// Struct to define interval time in updateTicker()
const (
	INTERVAL_PERIOD time.Duration = 24 * time.Hour
	HOUR_TO_TICK    int           = 00
	MINUTE_TO_TICK  int           = 00
	SECOND_TO_TICK  int           = 00
)

type Server struct {
	Tab_wd       *owm.All_data
	Lst_users    *users.All_users
	Lst_ball     *ballon.All_ball
	Lst_Devices  *devices.All_Devices
	Lst_workBall *ballonwork.All_work
	Logger       *log.Logger
}

/* Get the difference between Time.Now() et specifique time evenement and create a
** Tick channel of time package
 */
func updateTicker() *time.Ticker {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	diff := nextTick.Sub(time.Now())
	if diff <= 0 {
		diff = diff * -1
	}
	return time.NewTicker(diff)
}

/*
** Manage all go_routine Event.
 */
func Manage_goroutines(Serv *server.Server, base *db.Env) {
	channelfuncweatherdata := make(chan bool)
	channelfuncmoveball := make(chan bool)
	channelfuncupdatedata := make(chan bool)

	defer close(channelfuncmoveball)
	defer close(channelfuncweatherdata)
	defer close(channelfuncupdatedata)

	go func() {
		for {
			time.Sleep(2 * time.Minute)
			channelfuncmoveball <- true
		}
	}()
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			channelfuncweatherdata <- true
		}
	}()
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			channelfuncupdatedata <- true
		}
	}()
	go func() {
		ticker := updateTicker()
		for {
			<-ticker.C
			crontask.Send_AllBall(Serv)
			ticker = updateTicker()
		}
	}()

	for {
		select {
		case <-channelfuncweatherdata:
			{
				er := Serv.Tab_wd.Update_weather_data()
				if er != nil {
					Serv.Logger.Println("Update_weather error: ", er)
				} // If possible print Weather data with Serv.Tab_wd.Print_weatherdata()
				//				er = Serv.Lst_ball.Create_checkpoint(Serv.Tab_wd)
				//				if er != nil {
				//					Serv.Logger.Println("Create_checkpoint error: ", er)
				//				} // If possible print Checkpoint list with Serv.Lst_ball.Print_all_balls()
			}
		case <-channelfuncmoveball:
			{
				er := Serv.Lst_ball.Move_ball(Serv.Tab_wd)
				if er != nil {
					Serv.Logger.Println("Move_ball error: ", er)
				}
			}
		case <-channelfuncupdatedata:
			{
				fmt.Println("UPDATE !!")
				er := Serv.Lst_ball.Update_balls(Serv.Lst_ball, base)
				if er != nil {
					Serv.Logger.Println("Update_balls error: ", er)
				}
				er = Serv.Lst_users.Update_users(base)
				if er != nil {
					Serv.Logger.Println("Update_users error: ", er)
				}
			}
		}
	}
}

func ManageSignal(Serv *server.Server, myDb *db.Env) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	Serv.Logger.Println("Get Signal: ", s)

	fmt.Println("Save if possible begin crash")
	er := Serv.Lst_ball.Update_balls(Serv.Lst_ball, myDb)
	if er != nil {
		Serv.Logger.Println("Update_balls error: ", er)
	}
	er = Serv.Lst_users.Update_users(myDb)
	if er != nil {
		Serv.Logger.Println("Update_users error: ", er)
	}
	Serv.Tab_wd.Logger = nil
	Serv.Lst_users.Logger = nil
	Serv.Lst_ball.Logger = nil
	Serv.Lst_Devices.Logger = nil
	er = Serv.Lst_users.SaveUsersToFile()
	if er != nil {
		Serv.Logger.Println("Error on saveUsers, its a bad ! :(", er)
	}
	Serv.Logger = nil
	os.Exit(-1)
}

func main() {
	Server := new(server.Server)
	myDb := new(db.Env)

	er := Server.InitServer()
	if er != nil {
		fmt.Println("Error on InitServer: ", er)
		return
	}
	Db, er := myDb.OpenCo(er)
	if er != nil {
		Server.Logger.Println("OpenCo error: ", er)
		return
	}
	er = Server.Init_Data(myDb)
	if er != nil {
		Server.Logger.Println("Init_Data error: ", er)
		return
	}
	go ManageSignal(Server, myDb)
	go Manage_goroutines(Server, myDb)
	go sock.Listen(Server, Db)

	for {
		time.Sleep(time.Second * 60)
	}
}
