package server

import (
	"Wibo/ballon"
	"Wibo/ballonwork"
	"Wibo/db"
	"Wibo/debug"
	"Wibo/devices"
	"Wibo/owm"
	"Wibo/users"
	"container/list"
	"fmt"
	"log"
	"os"
)

type Server struct {
	Tab_wd       *owm.All_data
	Lst_users    *users.All_users
	Lst_ball     *ballon.All_ball
	Lst_Devices  *devices.All_Devices
	Lst_workBall *ballonwork.All_work
	Logger       *log.Logger
}

/*
** Initialize database before opening connection socket
 */
func (Serv *Server) Init_Data(base *db.Env) error {
	er := Serv.Tab_wd.Update_weather_data()
	if er != nil {
		return er
	} // If possible print Weather data with Serv.Tab_wd.Print_weatherdata()
	Serv.Tab_wd.Print_weatherdata()

	Serv.Lst_users.Ulist = list.New()
	er = Serv.Lst_users.Get_users(base.Db, Serv.Logger)
	if er != nil {
		return er
	} // If possible print user List with Serv.Lst_users.Print_users()
	Serv.Lst_users.GlobalStat = new(users.StatsUser)
	er = Serv.Lst_users.Get_GlobalStat(base)
	if er != nil {
		return er
	} // If possible print global stat with Serv.Lst_users.GlobalStat.Print()
	Serv.Lst_ball.Blist = list.New()
	er = Serv.Lst_ball.Get_balls(Serv.Lst_users, base)
	if er != nil {
		return er
	} // If possible print ball List with Serv.Lst_ball.Print_all_balls()
	Serv.Lst_Devices.Dlist = list.New()
	er = Serv.Lst_Devices.Get_devices(Serv.Lst_users, base)
	if er != nil {
		return er
	} // If possible print device List with Serv.Lst_Devices.Print_all_devices()
	Serv.Lst_workBall.Wlist = list.New()
	er = Serv.Lst_workBall.Get_workBall(base)
	if er != nil {
		return er
	} // If possible print workball List with Serv.Lst_Work.Print_all_workball()

	er = debug.CreateDataToDebug(Serv.Lst_ball, Serv.Lst_users, Serv.Tab_wd)
	if er != nil {
		return er
	} // If possible comment this section. Data filled to debug.
	er = Serv.Lst_ball.Create_checkpoint(Serv.Tab_wd)
	if er != nil {
		Serv.Logger.Println("Create_checkpoint error: ", er)
	} // If possible print Checkpoint list with Serv.Lst_ball.Print_all_balls()

	return nil
}

/*
** Prepare struct server to use it.
** Create a Loggerfile system.
 */
func (Serv *Server) InitServer() error {
	Serv.Tab_wd = new(owm.All_data)
	Serv.Lst_users = new(users.All_users)
	Serv.Lst_ball = new(ballon.All_ball)
	Serv.Lst_Devices = new(devices.All_Devices)
	Serv.Lst_workBall = new(ballonwork.All_work)

	file, er := os.Create("LogsSys.txt")
	if er != nil {
		fmt.Println("Erreur to create LogsSys.txt: ", er)
		return er
	}
	Serv.Logger = log.New(file, "logger: ", log.Lshortfile|log.Ldate|log.Ltime)
	file2, er := os.Create("LogsWeathers.txt")
	if er != nil {
		fmt.Println("Erreur to create LogsWeather.txt: ", er)
		return er
	}
	Serv.Tab_wd.Logger = log.New(file2, "Weather: ", log.Lshortfile|log.Ldate|log.Ltime)
	file3, er := os.Create("LogsBddBall.txt")
	if er != nil {
		fmt.Println("Erreur to create LogsBddBall.txt: ", er)
		return er
	}
	Serv.Lst_ball.Logger = log.New(file3, "BddBall: ", log.Lshortfile|log.Ldate|log.Ltime)
	file4, er := os.Create("LogsBddUsers.txt")
	if er != nil {
		fmt.Println("Erreur to create LogsBddUsers.txt: ", er)
		return er
	}
	Serv.Lst_users.Logger = log.New(file4, "BddUsers: ", log.Lshortfile|log.Ldate|log.Ltime)
	file5, er := os.Create("LogsBddDevice.txt")
	if er != nil {
		fmt.Println("Erreur to create LogsBddDevice.txt: ", er)
		return er
	}
	Serv.Lst_Devices.Logger = log.New(file5, "BddDevice: ", log.Lshortfile|log.Ldate|log.Ltime)
	Serv.Logger.Println("Init Server Done")
	return er
}
