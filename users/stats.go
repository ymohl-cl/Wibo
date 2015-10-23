package users

import (
"database/sql"
"time"
)

/*
Func Insert values into stats_users
Date YY/MM/DD
idusers_stats unique and foreing key table "users"

if  idusers_stats inserted once 
    else
error: Key (iduser_stats)=(28) already exists.


postgres sql
CREATE OR REPLACE FUNCTION public.setstatsuser(creationdate date, n_cont integer, n_cath integer, n_follow integer, n_message integer, n_send integer, iduser integer)
 RETURNS boolean
 LANGUAGE plpgsql
AS $function$
DECLARE
done boolean:= false;
BEGIN
 done :=  NOT exists(SELECT iduser_stats FROM stats_users WHERE iduser_stats=iduser);
IF done = false THEN
    UPDATE stats_users SET(num_owner, num_cath, num_follow, num_message, num_send) = ($2, $3, $4, $5, $6) WHERE iduser_stats=iduser;
    RETURN TRUE;
END IF;
PERFORM 1 FROM stats_users WHERE iduser_stats=iduser LIMIT 1;
IF NOT FOUND THEN 
INSERT INTO stats_users(creation_time, num_owner, num_cath, num_follow, num_message, num_send, iduser_stats) VALUES(creationdate, n_cont, n_cath, n_follow, n_message, n_send, iduser);
RETURN TRUE;
END IF;
RETURN FALSE;
END;
$function$
*/ 

func (Lusr *All_users) SetStatsByUser(c_idUser int64, u_stats *StatsUser, Db *sql.DB) bool {
	var err error
   stm , err := Db.Prepare("select setstatsuser($1, $2, $3, $4, $5, $6, $7)") 
    if (err != nil){
        return false
    }
  	_, err = stm.Query(u_stats.CreationDate, u_stats.NbrBallCreate, u_stats.NbrCatch, u_stats.NbrFollow, u_stats.NbrMessage, u_stats.NbrSend, c_idUser)
    
    if (err != nil){
        return false
    }
    return true
}   

func (Lusr *All_users) GetStatsByUser(idUser int64, Db *sql.DB) ( *StatsUser) {
	rows, err := Db.Query("SELECT creation_time, num_owner, num_catch, num_follow, num_message, num_send  FROM stats_users  WHERE iduser_stats = $1;", idUser)
	checkErr(err)
	for rows.Next(){
		var creationdate time.Time
		var ncontainers, ncath, nsend, nfollow, nmessage int64	
		err = rows.Scan(&creationdate,&ncontainers, &ncath, &nfollow, &nmessage, nsend)
		checkErr(err)
		return &StatsUser{CreationDate: creationdate, NbrBallCreate: ncontainers, NbrCatch: ncath, NbrSend: nsend, NbrFollow: nfollow, NbrMessage: nmessage}
	}
	return nil
}