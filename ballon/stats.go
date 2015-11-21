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
	defer stm.Close()
	row, err := stm.Query(b_stats.NbrKm, b_stats.NbrCatch, b_stats.NbrFollow, b_stats.NbrMagnet, c_idBall)
	if err != nil {
		return false
	}
	defer row.Close()
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
	defer rows.Close()
	if rows.Next() != false {
		rows.Scan(&postcreat, &datec)
		coord, err = GetCord(postcreat)
		if err != nil {
			return nil, time.Now(), err
		}
	}
	return coord, datec, nil
}

func (Lusr *All_ball) GetStatsBallon(idBall int64, Db *sql.DB) (*StatsBall, error) {
	fmt.Println("Id ball getStatBallon: ", idBall)
	rows, err := Db.Query("SELECT num_km, num_catch, num_follow, num_magnet  FROM stats_container  WHERE idball_stats=$1;", idBall)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var nkm float64
		var ncath, nfollow, nmagnet int64
		err = rows.Scan(&nkm, &ncath, &nfollow, &nmagnet)
		if err != nil {
			return nil, err
		}
		coord, date, err := GetCreationCoordDateBall(idBall, Db)
		if err != nil {
			return nil, err
		}
		fmt.Println("Coord created: ", coord)
		return &StatsBall{CreationDate: date, CoordCreated: coord, NbrKm: nkm, NbrFollow: nfollow, NbrCatch: ncath, NbrMagnet: nmagnet}, nil
	}
	return nil, err
}
