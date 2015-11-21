package users

import (
	"Wibo/db"
)

func (Lu *All_users) Get_GlobalStat(base *db.Env) error {
	rows, err := base.Db.Query("SELECT num_users, num_follow, num_message, num_send, num_cont, id_g, nbr_catch FROM globalstats ORDER BY id_g DESC LIMIT 1;")
	if err != nil {
		Lu.Logger.Println("err on Get_GlobalStat: ", err)
		return err
	}
	defer rows.Close()
	var idg int
	if rows.Next() {
		rows.Scan(&Lu.NbrUsers, &Lu.GlobalStat.NbrFollow, &Lu.GlobalStat.NbrMessage, &Lu.GlobalStat.NbrSend, &Lu.GlobalStat.NbrBallCreate, &idg, &Lu.GlobalStat.NbrCatch)
	}
	return err
}

// FUNCTION updatelocationuser(iduser integer, latitudec double precision, longitudec double precision)
// FUNCTION public.updateuser(iduser integer, latitudec double precision, longitudec double precision, log date)

func (lu *All_users) Update_users(base *db.Env) (err error) {
	u := lu.Ulist.Front()
	for u != nil {
		cu := u.Value.(*User)
		trow, err := base.Db.Query("SELECT updateuser($1, $2, $3, $4);", cu.Id, cu.Coord.Lon, cu.Coord.Lat, cu.Log)
		if err != nil {
			lu.Logger.Println("err Query updateUser: ", err)
			return err
		} else {
			defer trow.Close()
			ex := lu.SetStatsByUser(cu.Id, cu.Stats, base.Db)
			if ex != true {
				lu.Logger.Println("Fail to update user stats")
			}
		}
		u = u.Next()
	}
	lu.UpdateGlobal(base)
	return nil
}

func (lu *All_users) UpdateGlobal(base *db.Env) (err error) {
	trow, err := base.Db.Query("INSERT INTO globalstats(num_users, num_follow, num_message, num_send, num_cont, nbr_catch) VALUES($1, $2, $3, $4, $5, $6);",
		lu.Ulist.Len(),
		lu.GlobalStat.NbrFollow,
		lu.GlobalStat.NbrMessage,
		lu.GlobalStat.NbrSend,
		lu.GlobalStat.NbrBallCreate,
		lu.GlobalStat.NbrCatch)
	if err != nil {
		lu.Logger.Println("Update global stat wibo erreur: ", err)
		return err
	} else {
		defer trow.Close()
	}
	return nil
}
