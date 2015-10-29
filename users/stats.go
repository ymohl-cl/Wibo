package users

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
 "github.com/lib/pq"
_ "time"
"fmt"
// "log"
)

/*
Func Insert values into stats_users
Date YY/MM/DD
idusers_stats unique and foreing key table "users"

if  idusers_stats inserted once 
    else
error: Key (iduser_stats)=(28) already exists.


postgres sql
CREATE OR REPLACE FUNCTION public.setstatsuser(n_cont integer, n_catch integer, n_follow integer, n_message integer, n_send integer, iduser integer)
 RETURNS boolean
 LANGUAGE plpgsql
AS $function$
DECLARE
done boolean:= false;
BEGIN
 done :=  NOT exists(SELECT iduser_stats FROM stats_users WHERE iduser_stats=iduser);
IF done = false THEN
    UPDATE stats_users SET(num_owner, num_catch, num_follow, num_message, num_send) = ($1, $2, $3, $4, $5) WHERE iduser_stats=iduser;
    RETURN TRUE;
END IF;
PERFORM 1 FROM stats_users WHERE iduser_stats=iduser LIMIT 1;
IF NOT FOUND THEN 
INSERT INTO stats_users(num_owner, num_catch, num_follow, num_message, num_send, iduser_stats) VALUES(n_cont, n_catch, n_follow, n_message, n_send, iduser);
RETURN TRUE;
END IF;
RETURN FALSE;
END;
$function$
*/ 

func (Lusr *All_users) SetStatsByUser(c_idUser int64, u_stats *StatsUser, Db *sql.DB) bool {
	var err error
   stm , err := Db.Prepare("select setstatsuser($1, $2, $3, $4, $5, $6)") 
    if (err != nil){
        return false
    }
  	_, err = stm.Query(u_stats.NbrBallCreate, u_stats.NbrCatch, u_stats.NbrFollow, u_stats.NbrMessage, u_stats.NbrSend, c_idUser)
    
    if (err != nil){
        return false
    }
    return true
}

// CREATE FUNCTION create_statsballon() RETURNS trigger
//     LANGUAGE plpgsql
//     AS $$
//     BEGIN
//         --
//         -- Create a row in stats_container to reflect the operation performed on container,
//         -- make use of the special variable TG_OP to work out the operation.
//         --
       
//         IF (TG_OP = 'INSERT') THEN
//             INSERT INTO stats_container VAlUES (NEW.id, 0, 0, 0,false);
//             RETURN NEW;
//         END IF;
//         RETURN NULL; -- result is ignored since this is an AFTER trigger
//     END;
// $$;


func (Lusr *All_users) GetStatsByUser(idUser int64, Db *sql.DB) ( *StatsUser) {
    var creation pq.NullTime
    var ncontainers, ncath, nsend, nfollow, nmessage int
     _ = Db.QueryRow("SELECT  pg_catalog.date(\"user\".creationdate), num_owner, num_catch, num_follow, num_message, num_send  FROM stats_users INNER JOIN \"user\" ON  (stats_users.iduser_stats = \"user\".id_user) WHERE iduser_stats=$1;", int(idUser)).Scan(&creation, &ncontainers, &ncath, &nfollow, &nmessage, &nsend)
        fmt.Printf("\x1b[31;1m User: %v \x1b[0m\n", idUser)
        fmt.Printf("creationdate type: %T | Value:%v \n", creation, creation.Time)
        fmt.Printf("ncontainer type: %T | Value:%v \n", ncontainers, ncontainers)
        fmt.Printf("ncath type: %T | value: %v\n", ncath, ncath)
        fmt.Printf("nsend type: %T | value: %v \n", nsend, nsend)
        fmt.Printf("nfollow type: %T |  value: %v \n", nfollow, nfollow)
        fmt.Printf("nmeesage: type: %T | value: %v\n", nmessage, nmessage)
         return &StatsUser{CreationDate: creation.Time, NbrBallCreate: int64(ncontainers),
                 NbrCatch: int64(ncath), NbrSend: int64(nsend), NbrFollow: int64(nfollow),
                 NbrMessage: int64(nmessage)}
    return nil
    
}
