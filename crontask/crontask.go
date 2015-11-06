package crontask

import (
	"Wibo/ballon"
	"Wibo/server"
	"Wibo/users"
)

func Send_AllBall(Serv *server.Server) {
	lUser := Serv.Lst_users

	for el := lUser.Ulist.Front(); el != nil; el = el.Next() {
		user := el.Value.(*users.User)
		for eb := user.Possessed.Front(); eb != nil; eb = eb.Next() {
			ball := eb.Value.(*ballon.Ball)
			ball.Possessed = nil
			ball.InitCoord(user.Coord.Lon, user.Coord.Lat, 0, Serv.Tab_wd, true)
			user.Possessed.Remove(eb)
		}
		user.NbrBallSend = 0
	}
}
