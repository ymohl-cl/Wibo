package ballon

import (
	"database/sql"
)

/*
CREATE OR REPLACE FUNCTION public.setstatballon(n_km float, n_cath integer, n_follow integer, n_magnet integer, idball integer)
 RETURNS boolean
 LANGUAGE plpgsql
AS $function$
DECLARE
done boolean:= false;
BEGIN
 done :=  NOT exists(SELECT idball_stats FROM stats_container WHERE idball_stats=idball);
IF done = false THEN
    UPDATE stats_container SET(num_km, num_cath, num_follow, num_magnet) = ($1, $2, $3, $4) WHERE idball_stats=idball;
    RETURN TRUE;
END IF;
PERFORM 1 FROM stats_container WHERE idball_stats=idball LIMIT 1;
IF NOT FOUND THEN
INSERT INTO stats_container(num_km, num_cath, num_follow, num_magnet, idball_stats) VALUES(n_km, n_cath, n_follow, n_magnet, idball);
RETURN TRUE;
END IF;
RETURN FALSE;
END;
$function$
*/

func (Lusr *All_ball) SetStatsBallon(c_idBall int64, b_stats *StatsBall, Db *sql.DB) bool {
	var err error
	stm, err := Db.Prepare("select setstatballon($1, $2, $3, $4, $5)")
	if err != nil {
		return false
	}
	_, err = stm.Query(b_stats.NbrKm, b_stats.NbrCatch, b_stats.NbrFollow, b_stats.NbrMagnet, c_idBall)

	if err != nil {
		return false
	}
	return true
}

func (Lusr *All_ball) GetStatsBallon(idBall int64, Db *sql.DB) *StatsBall {
	rows, err := Db.Query("SELECT num_km, num_catch, num_follow, num_magnet  FROM stats_container  WHERE idball_stats = $1;", idBall)
	Lusr.checkErr(err)
	defer rows.Close()
	if rows.Next() != false {
		for rows.Next() {
			var nkm float64
			var ncath, nfollow, nmagnet int64
			err = rows.Scan(&nkm, &ncath, &nfollow, &nmagnet)
			Lusr.checkErr(err)
			return &StatsBall{NbrKm: nkm, NbrFollow: nfollow, NbrCatch: ncath, NbrMagnet: nmagnet}
		}
	}
	return nil
}
