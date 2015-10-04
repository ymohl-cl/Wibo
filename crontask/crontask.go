package crontask

import (
	"Wibo/ballon"
	"Wibo/owm"
	"Wibo/users"
)

func Send_AllBall(lBall *ballon.All_ball, lUser *users.All_users, Tab_wd *owm.All_data) {
	var checkpoint ballon.Checkpoint

	for el := lUser.Ulist.Front(); el != nil; el = el.Next() {
		user := el.Value.(*users.User)
		for eb := user.Possessed.Front(); eb != nil; eb = eb.Next() {
			ball := eb.Value.(*ballon.Ball)
			ball.Possessed = nil
			checkpoint.Coord.Lon = user.Coord.Lon
			checkpoint.Coord.Lat = user.Coord.Lat
			ball.Coord = ball.Checkpoints.PushBack(checkpoint)
			ball.Get_checkpointList(Tab_wd.Get_Paris())
			user.Possessed.Remove(eb)
		}
		user.NbrBallSend = 0
	}
}
