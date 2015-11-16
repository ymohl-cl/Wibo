package ballon

import (
	"database/sql"
	"fmt"
	"time"
)

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

func GetCreationCoordDateBall(idBall int64, Db *sql.DB) (*Coordinate, time.Time, error) {
	var coord *Coordinate
	var postcreat string
	var datec time.Time

	rows, err := Db.Query("SELECT ST_AsText(location_ckp), date FROM checkpoints WHERE checkpoints.id=(SELECT min(id) FROM checkpoints WHERE containerid=$1);", idBall)
	if err != nil {
		return nil, time.Now(), err
	}
	if rows.Next() != false {
		rows.Scan(&postcreat, &datec)
		coord = GetCord(postcreat)
	}
	return coord, datec, nil
}

func (Lusr *All_ball) GetStatsBallon(idBall int64, Db *sql.DB) *StatsBall {
	rows, err := Db.Query("SELECT num_km, num_cath, num_follow, num_magnet  FROM stats_container  WHERE idball_stats = $1;", idBall)
	Lusr.checkErr(err)
	defer rows.Close()
	if rows.Next() != false {
		for rows.Next() {
			var nkm float64
			var ncath, nfollow, nmagnet int64
			err = rows.Scan(&nkm, &ncath, &nfollow, &nmagnet)
			Lusr.checkErr(err)
			coord, date, err := GetCreationCoordDateBall(idBall, Db)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			return &StatsBall{CreationDate: date, CoordCreated: coord, NbrKm: nkm, NbrFollow: nfollow, NbrCatch: ncath, NbrMagnet: nmagnet}
		}
	}
	return nil
}
