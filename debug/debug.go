package debug

import (
	"Wibo/ballon"
	"Wibo/owm"
	"Wibo/users"
	"container/list"
	"time"
)

/* CREER UN BALLON POUR FAIRE DES TESTS */
func CreateDataToDebug(lball *ballon.All_ball, luser *users.All_users, tabwd *owm.All_data) error {
	/* CREATE USER */

	euser := luser.Ulist.Front()

	/* END CREATE USER -- CREATE CHECKPOINT */

	tmp_lst := list.New()
	var check_test0 ballon.Checkpoint
	var check_test1 ballon.Checkpoint
	var check_test2 ballon.Checkpoint
	var check_test3 ballon.Checkpoint
	var check_test4 ballon.Checkpoint
	var check_test5 ballon.Checkpoint

	check_test0.Coord.Lon = 2.316055
	check_test0.Coord.Lat = 48.833086
	check_test0.Date = time.Now()

	check_test1.Coord.Lon = 2.316065
	check_test1.Coord.Lat = 48.833586
	check_test1.Date = time.Now()

	check_test2.Coord.Lon = 2.30810777
	check_test2.Coord.Lat = 48.919253
	check_test2.Date = time.Now()

	check_test3.Coord.Lon = 2.3088211
	check_test3.Coord.Lat = 48.918361
	check_test3.Date = time.Now()

	check_test4.Coord.Lon = 2.316045
	check_test4.Coord.Lat = 48.833986
	check_test4.Date = time.Now()

	check_test5.Coord.Lon = 2.3080535
	check_test5.Coord.Lat = 48.910242
	check_test5.Date = time.Now()

	/* END CHECKPOINT -- BEGIN COORDINATE */

	var Coord0 ballon.Coordinate
	var Coord1 ballon.Coordinate
	var Coord2 ballon.Coordinate
	var Coord3 ballon.Coordinate
	var Coord4 ballon.Coordinate
	var Coord5 ballon.Coordinate

	Coord0.Lon = 2.3080535
	Coord0.Lat = 48.910242

	Coord1.Lon = 2.316055
	Coord1.Lat = 48.833086

	Coord2.Lon = 2.316065
	Coord2.Lat = 48.833586

	Coord3.Lon = 2.30810777
	Coord3.Lat = 48.919253

	Coord4.Lon = 2.3088211
	Coord4.Lat = 48.918361

	Coord5.Lon = 2.316045
	Coord5.Lat = 48.833986

	/* END COORDINATE -- BEGIN MESSAGE */

	mmp2 := list.New()
	mmp := list.New()
	var message0 ballon.Message
	var message1 ballon.Message
	var message2 ballon.Message
	var message3 ballon.Message

	message0.Id = 0
	message0.Size = 34
	message0.Content = "Coucou les gens, message0. ID 0..."
	message0.Type = 1

	message1.Id = 1
	message1.Size = 33
	message1.Content = "Coucou les gens, message1. Id 1.."
	message1.Type = 1

	message2.Id = 2
	message2.Size = 32
	message2.Content = "Coucou les gens, message2. Id 2."
	message2.Type = 1

	message3.Id = 3
	message3.Size = 31
	message3.Content = "Coucou les gens, message3. Id 3"
	message3.Type = 1

	mmp.PushBack(message0)
	mmp.PushBack(message3)
	mmp2.PushBack(message0)
	mmp2.PushBack(message1)
	mmp2.PushBack(message2)

	/* END MESSAGE -- BEGIN BALLON */

	ball0 := new(ballon.Ball)
	ball1 := new(ballon.Ball)
	ball2 := new(ballon.Ball)
	ball3 := new(ballon.Ball)
	ball4 := new(ballon.Ball)
	ball5 := new(ballon.Ball)

	ball0.Id_ball = 0
	ball0.Edited = false
	ball0.Title = "toto"
	ball0.Coord = tmp_lst.PushBack(check_test0)
	ball0.Wind = ballon.Wind{}
	ball0.Messages = mmp
	ball0.Date = time.Now()
	ball0.Checkpoints = list.New()
	ball0.Possessed = nil
	ball0.Followers = list.New()
	ball0.Creator = euser
	ball0.Stats = new(ballon.StatsBall)
	ball0.Stats.CreationDate = time.Now()
	ball0.Stats.CoordCreated = &Coord0

	ball1.Id_ball = 1
	ball1.Edited = false
	ball1.Title = "tata"
	ball1.Coord = tmp_lst.PushBack(check_test1)
	ball1.Wind = ballon.Wind{}
	ball1.Messages = mmp
	ball1.Date = time.Now()
	ball1.Checkpoints = list.New()
	ball1.Possessed = nil
	ball1.Followers = list.New()
	ball1.Creator = euser
	ball1.Stats = new(ballon.StatsBall)
	ball1.Stats.CreationDate = time.Now()
	ball1.Stats.CoordCreated = &Coord1

	ball2.Id_ball = 2
	ball2.Edited = false
	ball2.Title = "tutu"
	ball2.Coord = tmp_lst.PushBack(check_test2)
	ball2.Wind = ballon.Wind{}
	ball2.Messages = mmp
	ball2.Date = time.Now()
	ball2.Checkpoints = list.New()
	ball2.Possessed = nil
	ball2.Followers = list.New()
	ball2.Creator = euser
	ball2.Stats = new(ballon.StatsBall)
	ball2.Stats.CreationDate = time.Now()
	ball2.Stats.CoordCreated = &Coord2

	ball3.Id_ball = 3
	ball3.Edited = false
	ball3.Title = "tete"
	ball3.Coord = tmp_lst.PushBack(check_test3)
	ball3.Wind = ballon.Wind{}
	ball3.Messages = mmp
	ball3.Date = time.Now()
	ball3.Checkpoints = list.New()
	ball3.Possessed = nil
	ball3.Followers = list.New()
	ball3.Creator = euser
	ball3.Stats = new(ballon.StatsBall)
	ball3.Stats.CreationDate = time.Now()
	ball3.Stats.CoordCreated = &Coord3

	ball4.Id_ball = 4
	ball4.Edited = false
	ball4.Title = "tyty"
	ball4.Coord = tmp_lst.PushBack(check_test4)
	ball4.Wind = ballon.Wind{}
	ball4.Messages = mmp
	ball4.Date = time.Now()
	ball4.Checkpoints = list.New()
	ball4.Possessed = nil
	ball4.Followers = list.New()
	ball4.Creator = euser
	ball4.Stats = new(ballon.StatsBall)
	ball4.Stats.CreationDate = time.Now()
	ball4.Stats.CoordCreated = &Coord4

	ball5.Id_ball = 5
	ball5.Edited = false
	ball5.Title = "PROUT"
	ball5.Coord = tmp_lst.PushBack(check_test5)
	ball5.Wind = ballon.Wind{}
	ball5.Messages = mmp
	ball5.Date = time.Now()
	ball5.Checkpoints = list.New()
	ball5.Possessed = nil
	ball5.Followers = list.New()
	ball5.Creator = euser
	ball5.Stats = new(ballon.StatsBall)
	ball5.Stats.CreationDate = time.Now()
	ball5.Stats.CoordCreated = &Coord5

	// Add balls to list
	euser.Value.(*users.User).Followed.PushBack(lball.Blist.PushBack(ball0))
	euser.Value.(*users.User).Followed.PushBack(lball.Blist.PushBack(ball1))
	euser.Value.(*users.User).Followed.PushBack(lball.Blist.PushBack(ball2))
	euser.Value.(*users.User).Followed.PushBack(lball.Blist.PushBack(ball3))
	euser.Value.(*users.User).Followed.PushBack(lball.Blist.PushBack(ball4))
	euser.Value.(*users.User).Followed.PushBack(lball.Blist.PushBack(ball5))
	euser.Value.(*users.User).Stats.NbrBallCreate = 5
	euser.Value.(*users.User).Stats.NbrSend = 5
	euser.Value.(*users.User).Stats.NbrMessage = 5

	/* END BALL -- BEGIN CREATE_CHECKPOINT */

	er := lball.Create_checkpoint(tabwd)
	// If possible print ball List with Serv.Lst_ball.Print_all_balls()
	return er
}
