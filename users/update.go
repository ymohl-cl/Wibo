package users

import (
	"Wibo/db"
	"fmt"
)

func (Lu *All_users) Get_GlobalStat(base *db.Env) error {
	rows, err := base.Db.Query("SELECT num_users, num_follow, num_message, num_send, num_cont FROM globalStats ORDER BY lastupdate DESC LIMIT 1;")
	if err != nil {
		return &userError{Prob: "Get Global stat", Err: nil, Logf: Lu.Logger}
	}
	defer rows.Close()
	rows.Scan(&Lu.NbrUsers, &Lu.GlobalStat.NbrFollow, &Lu.GlobalStat.NbrMessage, &Lu.GlobalStat.NbrSend, &Lu.GlobalStat.NbrBallCreate)
	return nil
}

// FUNCTION updatelocationuser(iduser integer, latitudec double precision, longitudec double precision)
// FUNCTION public.updateuser(iduser integer, latitudec double precision, longitudec double precision, log date)

func (lu *All_users) Update_users(base *db.Env) (err error) {
	u := lu.Ulist.Front()
	for u != nil {
		cu := u.Value.(*User)
		trow, err := base.Db.Query("SELECT updateuser($1, $2, $3, $4);", cu.Id, cu.Coord.Lon, cu.Coord.Lat, cu.Log)
		if err != nil {
			return &userError{Prob: "Update users", Err: err, Logf: lu.Logger}
		}

		defer trow.Close()
		ex := lu.SetStatsByUser(cu.Id, cu.Stats, base.Db)
		if ex != true {
			fmt.Println("Fail to update user stats")
		}
		u = u.Next()
	}
	lu.UpdateGlobal(base)
	return nil
}

func (lu *All_users) UpdateGlobal(base *db.Env) (err error) {
	trow, err := base.Db.Query("INSERT INTO globalstats(num_users, num_follow, num_mesage, num_send, num_cont) VALUES($1, $2, $3, $4, $5);",
		lu.NbrUsers,
		lu.GlobalStat.NbrFollow,
		lu.GlobalStat.NbrMessage,
		lu.GlobalStat.NbrSend,
		lu.GlobalStat.NbrBallCreate)
	if err != nil {
		return &userError{Prob: "Update global wibo", Err: err, Logf: lu.Logger}
	}
	defer trow.Close()
	return nil
}
