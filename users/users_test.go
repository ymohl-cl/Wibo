package users_test

import (
	"Wibo/ballon"
	"Wibo/db"
	"Wibo/users"
	"container/list"
	"fmt"
	"testing"
	"time"
)

func TestUsers(t *testing.T) {

	var err error
	Lst_users := new(users.All_users)
	Lst_ball := new(ballon.All_ball)
	myDb := new(db.Env)
	Lst_users.Ulist = list.New()
	Lst_ball.Blist = list.New()
	Db, err := myDb.OpenCo(err)
	user1 := new(users.User)
	user1.Mail = "Toto2@Dr.fr"
	user1.Log = time.Now()
	user1.Followed = list.New()
	user1.Stats = &users.StatsUser{}
	user1.Stats.CreationDate = time.Now()
	user1.Coord.Lat = 48.833086
	user1.Coord.Lon = 2.316055
	user1.Log = time.Now()
	Lst_users.Ulist.PushBack(user1)

	/* CREER UN BALLON POUR FAIRE DES TESTS */

	result, _ := Lst_users.Add_new_user(user1, Db)
	fmt.Println(result)
}
