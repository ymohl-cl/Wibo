package devices_test

import (
	"Wibo/db"
	"Wibo/users"
	"container/list"
	// "time"
	"Wibo/ballon"
	"fmt"
	"testing"
	"Wibo/protocol"
	"Wibo/devices"
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
	Lprotocol := new(protocol.Request)
	Ldevices := new(devices.All_Devices)
	myDb := new(db.Env)
	Lst_users.Ulist = list.New()
	Lst_ball.Blist = list.New()
	Db, err := myDb.OpenCo(err)
	Lst_ball.Get_balls(Lst_users, myDb)
	fmt.Println(Db)

	Lprotocol.Coord.Lon = 48.833986
	Lprotocol.Coord.Lat = 2.316045
// user1 := new(users.User)
// 	user1.Id = 68
// 	user1.Mail = "mailtest2@test.com"
// 	user1.Log = time.Now()
// 	user1.Followed = list.New()
// user2 := new(users.User)
// 	user1.Id = 68
// 	user1.Mail = "mailtest3@test.com"
// 	user1.Log = time.Now()
// 	user1.Followed = list.New()
// user3 := new(users.User)
// 	user1.Id = 68
// 	user1.Mail = "mailtest4@test.com"
// 	user1.Log = time.Now()
// 	user1.Followed = list.New()
// user4 := new(users.User)
// 	user1.Id = 68
// 	user1.Mail = "mailtest5@test.com"
// 	user1.Log = time.Now()
// 	user1.Followed = list.New()
// user5 := new(users.User)
// 	user1.Id = 68
// 	user1.Mail = "mailtest6@test.com"
// 	user1.Log = time.Now()
// 	user1.Followed = list.New()
// user6 := new(users.User)
// 	user1.Id = 68
// 	user1.Mail = "mailtest7@test.com"
// 	user1.Log = time.Now()
// 	user1.Followed = list.New()
// user7 := new(users.User)
// 	user1.Id = 68
// 	user1.Mail = "mailtest8@test.com"
// 	user1.Log = time.Now()
// 	user1.Followed = list.New()

// 	Lst_users.Ulist.PushBack(user1)
// 	Lst_users.Ulist.PushBack(user1)
// 	Lst_users.Ulist.PushBack(user1)
// 	Lst_users.Ulist.PushBack(user1)
// 	Lst_users.Ulist.PushBack(user1)
	Lst_users.Get_users(myDb.Db)
	Ldevices.Get_devices(Lst_users, myDb)


}