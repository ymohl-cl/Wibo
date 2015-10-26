package users

import (
"database/sql"
"time"
"strings"
"fmt"
"log"
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


func GetDateFormat(qdate string) (fdate time.Time) {
    f := func(c rune) bool {
        return c == '"'
    }
    fields := strings.FieldsFunc(qdate, f)
    for _, value := range fields {
        qdate = string(value)
    }
    fdate, err := time.Parse("2006-01-02 15:04:05", qdate)
    checkErr(err)
    return fdate
}

func (Lusr *All_users) GetStatsByUser(idUser int64, Db *sql.DB) ( *StatsUser) {
    var creation time.Time
    var ncontainers, ncath, nsend, nfollow, nmessage int64
	err := Db.QueryRow("SELECT \"user\".creationdate, num_owner, num_catch, num_follow, num_message, num_send  FROM stats_users INNER JOIN \"user\" ON  (stats_users.iduser_stats = \"user\".id_user) WHERE iduser_stats=$1;", idUser).Scan(&creation, &ncontainers, &ncath, &nfollow, &nmessage, &nsend)
    switch {
    case err == sql.ErrNoRows:
        fmt.Println(err)
        return nil // make an error here
		      // return &StatsUser{NbrBallCreate: 0, NbrCatch:0, NbrSend: 0, NbrFollow: 0, NbrMessage: 0}
    case err != nil:
            log.Fatal(err)
    default:
            fmt.Printf("Username is %s | %s | %s | %s | %s | %v \n", creation, ncontainers,ncath, nsend, nfollow, nmessage)
        return &StatsUser{CreationDate: creation, NbrBallCreate: ncontainers, NbrCatch: ncath, NbrSend: nsend, NbrFollow: nfollow, NbrMessage: nmessage}
    }
    return nil
    
}
