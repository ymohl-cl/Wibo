package users

import (
	"database/sql"
	"github.com/lib/pq"
)

func (Lusr *All_users) SetStatsByUser(c_idUser int64, u_stats *StatsUser, Db *sql.DB) bool {
	var err error
	_, err = Db.Exec("select setstatsuser($1, $2, $3, $4, $5, $6);", u_stats.NbrBallCreate, u_stats.NbrCatch, u_stats.NbrFollow, u_stats.NbrMessage, u_stats.NbrSend, c_idUser)
	if err != nil {
		Lusr.Logger.Println(err)
		return false
	}
	return true
}

func (Lusr *All_users) GetStatsByUser(idUser int64, Db *sql.DB) *StatsUser {
	var creation pq.NullTime
	var ncontainers, ncath, nsend, nfollow, nmessage int
	err := Db.QueryRow("SELECT  pg_catalog.date(\"user\".creationdate), num_owner, num_catch, num_follow, num_message, num_send  FROM stats_users INNER JOIN \"user\" ON  (stats_users.iduser_stats = \"user\".id_user) WHERE iduser_stats=$1;", int(idUser)).Scan(&creation,
		&ncontainers,
		&ncath,
		&nfollow,
		&nmessage,
		&nsend)
	if err != nil {
		Lusr.Logger.Println("Error Query on GetStatsByUser: ", err)
	}
	return &StatsUser{CreationDate: creation.Time, NbrBallCreate: int64(ncontainers),
		NbrCatch: int64(ncath), NbrSend: int64(nsend), NbrFollow: int64(nfollow),
		NbrMessage: int64(nmessage)}
}
