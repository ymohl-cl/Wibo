package users

import (
	"Wibo/db"
	"fmt"
)

func (Lu *All_users) Get_GlobalStat(base *db.Env) error {
	rows, err := base.Db.Query("SELECT num_users, num_follow, num_message, num_send, num_cont, id_g FROM globalstats ORDER BY id_g DESC LIMIT 1;")
	if err != nil {
		Lu.Logger.Println("err on Get_GlobalStat: ", err)
		return err
		//return &userError{Prob: "Get Global stat", Err: nil, Logf: Lu.Logger}
	}
	defer rows.Close()
	var idg int
	if rows.Next() {
		rows.Scan(&Lu.NbrUsers, &Lu.GlobalStat.NbrFollow, &Lu.GlobalStat.NbrMessage, &Lu.GlobalStat.NbrSend, &Lu.GlobalStat.NbrBallCreate, &idg)
	}
	fmt.Println("Get global; statidg: ", idg)
	fmt.Println("NbrUser: ", Lu.NbrUsers)
	fmt.Println("NbrFollow", Lu.GlobalStat.NbrFollow)
	fmt.Println("NbrMessage", Lu.GlobalStat.NbrMessage)
	fmt.Println("NbrSend", Lu.GlobalStat.NbrSend)
	fmt.Println("NbrBallCreate", Lu.GlobalStat.NbrBallCreate)
	Lu.Logger.Println("Get stat ok: ", Lu.GlobalStat)
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
	fmt.Println("Insert global; stat")
	fmt.Println("NbrUser: ", lu.NbrUsers)
	fmt.Println("NbrUser: ", lu.Ulist.Len())
	fmt.Println("NbrFollow", lu.GlobalStat.NbrFollow)
	fmt.Println("NbrMessage", lu.GlobalStat.NbrMessage)
	fmt.Println("NbrSend", lu.GlobalStat.NbrSend)
	fmt.Println("NbrBallCreate", lu.GlobalStat.NbrBallCreate)
	trow, err := base.Db.Query("INSERT INTO globalstats(num_users, num_follow, num_message, num_send, num_cont) VALUES($1, $2, $3, $4, $5);",
		lu.Ulist.Len(),
		lu.GlobalStat.NbrFollow,
		lu.GlobalStat.NbrMessage,
		lu.GlobalStat.NbrSend,
		lu.GlobalStat.NbrBallCreate)
	if err != nil {
		lu.Logger.Println("Update global stat wibo erreur: ", err)
		return err
		//return &userError{Prob: "Update global wibo", Err: err, Logf: lu.Logger}
	}
	defer trow.Close()
	return nil
}
