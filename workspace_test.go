package main_test

import (
	"ballon"
	"container/list"
	//"owm"
	"testing"
	"time"
	//"users"
)

func TestBallon(t *testing.T) {
	//	Tab_wd := new(owm.All_data)
	Lst_ball := new(ballon.All_ball)
	//	Lst_users := new(users.All_users)
	/* CREER UN BALLON POUR FAIRE DES TESTS */
	tmp_lst := list.New()
	var check_test0 ballon.Checkpoint
	check_test0.Coord.Lon = 48.833086
	check_test0.Coord.Lat = 2.316055
	check_test0.Date = time.Now()
	/*var check_test1 ballon.Checkpoint
	check_test1.Coord.Lon = 48.833586
	check_test1.Coord.Lat = 2.316065
	check_test1.Date = time.Now()
	var check_test2 ballon.Checkpoint
	check_test2.Coord.Lon = 48.833368
	check_test2.Coord.Lat = 2.316059
	check_test2.Date = time.Now()
	var check_test3 ballon.Checkpoint
	check_test3.Coord.Lon = 48.833286
	check_test3.Coord.Lat = 2.316903
	check_test3.Date = time.Now()
	var check_test4 ballon.Checkpoint
	check_test4.Coord.Lon = 48.833986
	check_test4.Coord.Lat = 2.316045
	check_test4.Date = time.Now()
	mmp2 := list.New()
	mmp := list.New()
	var message0 ballon.Message
	message0.Id = 0
	message0.Size = 68
	message0.Content = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean com"
	message0.Type = 1
	var message1 ballon.Message
	message1.Id = 0
	message1.Size = 68
	message1.Content = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean com"
	message1.Type = 1
	var message2 ballon.Message
	message2.Id = 2
	message2.Size = 68
	message2.Content = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean com"
	message2.Type = 1

	mmp.PushBack(message0)
	mmp2.PushBack(message0)
	mmp2.PushBack(message1)
	mmp2.PushBack(message2)
	*/
	ball0 := new(ballon.Ball)
	ball0.Id_ball = 0
	ball0.Title = "toto"
	ball0.Coord = tmp_lst.PushBack(check_test0)
	ball0.Wind = ballon.Wind{}
	ball0.Messages = nil
	ball0.Date = time.Now()
	ball0.Checkpoints = list.New()
	ball0.Possessed = nil
	ball0.Followers = list.New()
	ball0.Creator = nil
	Lst_ball.Blist.PushBack(ball0)

	/* FIN DE LA CREATION DEBALLON POUR TEST */

}
