package users_test

import (
	"Wibo/ballon"
	"Wibo/db"
	"Wibo/users"
	"container/list"
	"fmt"
	"testing"
)

// CREATE OR REPLACE FUNCTION create_statsuser() RETURNS TRIGGER AS $stats_creation$
//     BEGIN
//         --
//         -- Create a row in stats_users to reflect the operation performed on user,
//         -- make use of the special variable TG_OP to work out the operation.
//         --

//         IF (TG_OP = 'INSERT') THEN
//             INSERT INTO stats_users VAlUES (0, 0, 0, 0, NEW.id_user);
//             RETURN NEW;
//         END IF;
//         RETURN NULL; -- result is ignored since this is an AFTER trigger
//     END;
// $stats_creation$ LANGUAGE plpgsql;

func TestUsers(t *testing.T) {

	var err error
	Lst_users := new(users.All_users)
	Lst_ball := new(ballon.All_ball)
	myDb := new(db.Env)
	Lst_users.Ulist = list.New()
	Lst_ball.Blist = list.New()
	Db, err := myDb.OpenCo(err)
	user1 := new(users.User)
	user1.Id = 125
	user1.Mail = "Toto@Dr.fr"
	user1.Log = time.Now()
	user1.Followed = list.New()
	user1.Stats = &users.StatsUser{}
	user1.Stats.CreationDate = time.Now()
	user1.Coord.Lat = 48.833086
	user1.Coord.Lon = 2.316055
	user1.Log = time.Now()
	Lst_users.Ulist.PushBack(user1)

	result := Lst_users.Add_new_user(user1, Db)
	fmt.Println(result)
}
