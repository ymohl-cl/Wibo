package crontask

import (
	"Wibo/ballon"
	"Wibo/server"
	"Wibo/users"
)

func Send_AllBall(Serv *server.Server) {
	var checkpoint ballon.Checkpoint
	lUser := Serv.Lst_users

	for el := lUser.Ulist.Front(); el != nil; el = el.Next() {
		user := el.Value.(*users.User)
		for eb := user.Possessed.Front(); eb != nil; eb = eb.Next() {
			ball := eb.Value.(*ballon.Ball)
			ball.Possessed = nil
			checkpoint.Coord.Lon = user.Coord.Lon
			checkpoint.Coord.Lat = user.Coord.Lat
			ball.Coord = ball.Checkpoints.PushBack(checkpoint)
			ball.Get_checkpointList(Serv.Tab_wd.Get_Paris())
			user.Possessed.Remove(eb)
		}
		user.NbrBallSend = 0
	}
}
