package users

import (
"database/sql"
"time"
)

/*

under construction please don't touch


CREATE FUNCTION update_stats() RETURNS trigger AS $update_stats$
    BEGIN
        -- Check that creationdate is not null
        IF NEW.creationdate IS NULL THEN
            RAISE EXCEPTION 'creationdate should not be null';
        END IF;
        -- Check index fourni pour les serveur
        IF NEW.ianix IS NULL THEN
            RAISE EXCEPTION 'Ballon % should have a index', NEW.titlename;
        END IF;

        IF NEW.ianix  < 0 THEN
            RAISE EXCEPTION 'Ballon % cannot have and index negative', NEW.titlename;
        END IF;

        -- Check title
        IF NEW.titlename IS NULL THEN
            RAISE EXCEPTION 'Ballon with % should have a title', NEW.ianix;
        END IF;

        RETURN NEW;
    END;
$emp_stamp$ LANGUAGE plpgsql;
*/

func (Lusr *All_users) GetStatsByUser(idUser int64, Db *sql.DB) ( *StatsUser) {
	rows, err := Db.Query("SELECT creationdate, n_containers, n_catch, n_send, n_follow, n_meesage  FROM statsuser;")
	checkErr(err)
	for rows.Next(){
		var creationdate time.Time
		var ncontainers, ncath, nsend, nfollow, nmessage int64	
		err = rows.Scan(&creationdate,&ncontainers, &ncath, &nsend, &nfollow, &nmessage)
		checkErr(err)
		return &StatsUser{CreationDate: creationdate, NbrBallCreate: ncontainers, NbrCatch: ncath, NbrSend: nsend, NbrFollow: nfollow, NbrMessage: nmessage}
	}
	return nil
}