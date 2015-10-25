package ballon

import (
"database/sql"
)

/*
CREATE OR REPLACE FUNCTION public.setstatballon(n_km integer, n_cath integer, n_follow integer, n_message integer, n_send integer, idball integer)
 RETURNS boolean
 LANGUAGE plpgsql
AS $function$
DECLARE
done boolean:= false;
BEGIN
 done :=  NOT exists(SELECT idball_stats FROM stats_container WHERE idball_stats=idball);
IF done = false THEN
    UPDATE stats_container SET(num_km, num_cath, num_follow, num_message, num_send) = ($1, $2, $3, $4, $5) WHERE idball_stats=idball;
    RETURN TRUE;
END IF;
PERFORM 1 FROM stats_container WHERE idball_stats=idball LIMIT 1;
IF NOT FOUND THEN 
INSERT INTO stats_container(num_km, num_cath, num_follow, num_message, num_send, idball_stats) VALUES(n_km, n_cath, n_follow, n_message, n_send, idball);
RETURN TRUE;
END IF;
RETURN FALSE;
END;
$function$
*/

func (Lusr *All_ball) SetStatsBallon(c_idBall int64, b_stats *StatsBall, Db *sql.DB) bool {
	var err error
	stm , err := Db.Prepare("select setstatballon($1, $2, $3, $4, $5, $6)") 
    if (err != nil){
        return false
    }
  	_, err = stm.Query(b_stats.NbrKm, b_stats.NbrCatch, b_stats.NbrFollowers, b_stats.NbrMessage, b_stats.NbrSend, c_idBall)
    
    if (err != nil){
        return false
    }
    return true
}   

func (Lusr *All_ball) GetStatsBallon(idBall int64, Db *sql.DB) ( *StatsBall) {
	rows, err := Db.Query("SELECT num_km, num_catch, num_follow, num_message, num_send  FROM stats_container  WHERE idball_stats = $1;", idBall)
	checkErr(err)
	for rows.Next(){
		var nkm, ncath, nsend, nfollow, nmessage int64	
		err = rows.Scan(&nkm, &ncath, &nfollow, &nmessage, nsend)
		checkErr(err)
		return &StatsBall{NbrKm: nkm, NbrFollowers: nfollow, NbrSend: nsend, NbrMessage: nmessage, NbrCatch: ncath}
	}
	return nil
}