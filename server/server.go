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
	Serv.Lst_users.Ulist = list.New()
	er = Serv.Lst_users.Get_users(base.Db)
	if er != nil {
		return er
	} // If possible print user List with Serv.Lst_users.Print_users()
	Serv.Lst_users.GlobalStat = new(users.StatsUser)
	er = Serv.Lst_users.Get_GlobalStat(base)
	if er != nil {
		return er
	} // If possible print global stat with Serv.Lst_users.GlobalStat.Print()
	Serv.Lst_ball.Blist = list.New()
	//	er = Serv.Lst_ball.Get_balls(Serv.Lst_users, base)
	//	if er != nil {
	//		return er
	//	} // If possible print ball List with Serv.Lst_ball.Print_all_balls()
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
		fmt.Println("Erreur to create LogsSys.txt")
		fmt.Println(er)
		return er
	}
	Serv.Logger = log.New(file, "logger: ", log.Llongfile|log.Ldate|log.Ltime)
	return er
}
