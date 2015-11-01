package ballon_test

import (
	"Wibo/ballon"
	"Wibo/db"
	"Wibo/users"
	"container/list"
	"fmt"
	_ "github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

// request data base test
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestBallon(t *testing.T) {

	var err error

	Lst_users := new(users.All_users)
	Lst_ball := new(ballon.All_ball)
	myDb := new(db.Env)
	Lst_users.Ulist = list.New()
	Lst_ball.Blist = list.New()
	Db, err := myDb.OpenCo(err)
	Lst_ball.Get_balls(Lst_users, myDb)
	fmt.Println(Db)

	user1 := new(users.User)
	user1.Id = 88
	user1.Mail = "mailtest4@test.com"
	user1.Log = time.Now()
	user1.Followed = list.New()
	user1.Stats = &users.StatsUser{}
	user1.Stats.CreationDate = time.Now()
	user1.Coord.Lat = 48.833086
	user1.Coord.Lon = 2.316055
	user1.Log = time.Now()
	pass := []byte("Pass1Test")
	fmt.Println("password test")
	fmt.Println(pass)
	bpass, err := bcrypt.GenerateFromPassword(pass, 15)
	/*
		fmt.Println(len(bpass))
		if err != nil {
			t.Fatalf("GenerateFromPassword error: %s", err)
		}
		if bcrypt.CompareHashAndPassword(bpass, pass) != nil {
			t.Errorf("%v should hash %s correctly", bpass, pass)
		}
		Lst_users.Ulist.PushBack(user1)
		b, err := Lst_users.Add_new_user(user1, Db, "Pass1Test")
		if err != nil {
			t.Fatalf("add user fail error: %s", err)
		}
		fmt.Println(b)
	*/
	rows, err := Db.Query("SELECT id_user, login, mail, passbyte FROM \"user\" WHERE id_user=$1;", 94)
	for rows.Next() {
		var idUser int64
		var login string
		var mailq string
		var passbyte []byte
		err = rows.Scan(&idUser, &login, &mailq, &passbyte)
		if bcrypt.CompareHashAndPassword(bpass, passbyte) != nil {
			t.Errorf("%v should hash %s correctly", bpass, passbyte)
		}
		fmt.Printf(" %v | %v", bpass, passbyte)
	}

	/* CREER UN BALLON POUR FAIRE DES TESTS */
	// tmp_lst := list.New()
	// var check_test0 ballon.Checkpoint
	// check_test0.Coord.Lon = 48.833086
	// check_test0.Coord.Lat = 2.316055
	// check_test0.Date = time.Now()

	// var check_test1 ballon.Checkpoint
	// check_test1.Coord.Lon = 48.833586
	// check_test1.Coord.Lat = 2.316065
	// check_test1.Date = time.Now()

	// var check_test2 ballon.Checkpoint
	// check_test2.Coord.Lon = 48.833368
	// check_test2.Coord.Lat = 2.316059
	// check_test2.Date = time.Now()

	// var check_test3 ballon.Checkpoint
	// check_test3.Coord.Lon = 48.833286
	// check_test3.Coord.Lat = 2.316903
	// check_test3.Date = time.Now()

	// var check_test4 ballon.Checkpoint
	// check_test4.Coord.Lon = 48.833986
	// check_test4.Coord.Lat = 2.316045
	// check_test4.Date = time.Now()

	// lmessages := list.New()
	// listMessage1 := list.New()
	// var message0 ballon.Message
	// message0.Id = 0
	// message0.Size = 68
	// message0.Content = "Mensaje 1 test"
	// message0.Type = 1
	// var message1 ballon.Message
	// message1.Id = 0
	// message1.Size = 68
	// message1.Content = "Mensaje 2 test"
	// message1.Type = 1
	// var message2 ballon.Message
	// message2.Id = 2
	// message2.Size = 68
	// message2.Content = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean com"
	// message2.Type = 1

	// listMessage1.PushBack(message0)
	// lmessages.PushBack(message0)
	// lmessages.PushBack(message1)
	// lmessages.PushBack(message2)

	// ball0 := new(ballon.Ball)
	// ball0.Id_ball = 45
	// ball0.Title = "toto"
	// ball0.Coord = tmp_lst.PushBack(check_test0)
	// ball0.Wind = ballon.Wind{}
	// ball0.Messages = listMessage1
	// ball0.Date = time.Now()
	// ball0.Checkpoints = list.New()
	// ball0.Possessed = nil
	// ball0.Followers = list.New()
	// ball0.Creator = nil
	// Lst_ball.Blist.PushBack(ball0)

	// ball1 := new(ballon.Ball)
	// ball1.Id_ball = 1
	// ball1.Title = "tata"
	// ball1.Coord = tmp_lst.PushBack(check_test1)
	// ball1.Wind = ballon.Wind{}
	// ball1.Messages = listMessage1
	// ball1.Date = time.Now()
	// ball1.Checkpoints = list.New()
	// ball1.Possessed = nil
	// ball1.Followers = list.New()
	// ball1.Creator = nil
	// Lst_ball.Blist.PushBack(ball1)

	// ball2 := new(ballon.Ball)
	// ball2.Id_ball = 2
	// ball2.Title = "tutu"
	// ball2.Coord = tmp_lst.PushBack(check_test2)
	// ball2.Wind = ballon.Wind{}
	// ball2.Messages = listMessage1
	// ball2.Date = time.Now()
	// ball2.Checkpoints = list.New()
	// ball2.Possessed = nil
	// ball2.Followers = list.New()
	// ball2.Creator = nil
	// Lst_ball.Blist.PushBack(ball2)

	// ball3 := new(ballon.Ball)
	// ball3.Id_ball = 3
	// ball3.Title = "tete"
	// ball3.Coord = tmp_lst.PushBack(check_test3)
	// ball3.Wind = ballon.Wind{}
	// ball3.Messages = listMessage1
	// ball3.Date = time.Now()
	// ball3.Checkpoints = list.New()
	// ball3.Possessed = nil
	// ball3.Followers = list.New()
	// ball3.Creator = nil
	// Lst_ball.Blist.PushBack(ball3)

	// Ball4 := new(ballon.Ball)
	// Ball4.Id_ball = 63
	// Ball4.Title = "tyty"
	// Ball4.Edited = true
	// Ball4.Coord = tmp_lst.PushBack(check_test4)
	// Ball4.Wind = ballon.Wind{Speed: 32, Degress: 2}
	// Ball4.Messages = listMessage1
	// Ball4.Date = time.Now()
	// Ball4.Checkpoints = tmp_lst
	// Ball4.Possessed = Lst_users.Ulist.Front()
	// Ball4.Followers = list.New()
	// Ball4.Creator = Lst_users.Ulist.Front()
	// Lst_ball.Blist.PushBack(Ball4)
	// if _, err := Lst_ball.InsertBallon(Ball4, myDb); err != nil {
	// 	t.Fatalf("Fail insert ball:%s", err)
	// }

	// if _, err :=	ball4.GetItinerary(myDb.Db); err != nil {
	// 	t.Fatalf("Fail get GetItinerary error: %s", err)
	// }
	// Lst_ball.Update_balls(Lst_ball, myDb)
	// fmt.Println("\x1b[31;1m SECOND PRINT ALL BALLS\x1b[0m")
	//Lst_ball.Print_all_balls()

}

// func BenchmarkAddNewDefaultUser(b *testing.B){
// 		var err error
// 		Lst_users := new(users.All_users)
// 		myDb := new(db.Env)
// 		Lst_users.Ulist = list.New()
// 		Db, err := myDb.OpenCo(err)
// 		if err != nil {
// 			b.Fatalf("benchmarkConnection: %s", err)
// 		}
// 		  for n := 0; n < b.N; n++ {
// 				defU := Lst_users.AddNewDefaultUser(Db)
// 				// if err != nil {
// 				// b.Fatalf("benchmarkAddNewDefaultUser: %s", err)
// 				// }S
// 			fmt.Println(defU)
//         }
// }

// func checkUpdate_balls(t *testing.T, LBall *ballon.All_ball, base *db.Env){
// 		LBall.Update_balls(LBall, base)
// }
