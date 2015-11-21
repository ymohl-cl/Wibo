package ballon

import (
	"Wibo/db"
	"Wibo/users"
	"container/list"
	"database/sql"
	_ "fmt"
	"strings"
)

func (Lst_ball *All_ball) InsertMessages(messages *list.List, idBall int64, base *db.Env) (err error) {
	i := 0
	for e := messages.Front(); e != nil; e = e.Next() {
		err = base.Transact(base.Db, func(tx *sql.Tx) error {
			stm, err := tx.Prepare("INSERT INTO message(content, containerid, index_m, size) VALUES ($1, $2, $3, $4)")
			if err != nil {
				Lst_ball.Logger.Println("Erreur tx prepare: ", err)
				return err
			}
			defer stm.Close()
			row, err := stm.Query(strings.Trim(e.Value.(Message).Content, "\x00"), idBall, i, e.Value.(Message).Size)
			if err != nil {
				Lst_ball.Logger.Println("Erreur Query: ", err)
				return err
			}
			defer row.Close()
			i++
			return err
		})
	}
	return nil
}

func (Lstb *All_ball) SetFollowerBalls(curr_b *Ball, base *db.Env) {
	var idB int64
	base.Db.QueryRow("SELECT id FROM container WHERE ianix=$1;", curr_b.Id_ball).Scan(&idB)
	row, er := base.Db.Query("DELETE FROM followed WHERE container_id=$1", idB)
	if er != nil {
		Lstb.Logger.Println("Error on Delete followed: ", er)
	} else {
		defer row.Close()
	}
	for f := curr_b.Followers.Front(); f != nil; f = f.Next() {
		err := base.Transact(base.Db, func(tx *sql.Tx) error {
			stm, err := tx.Prepare("INSERT INTO followed(container_id, iduser) values($1,$2)")
			if err != nil {
				Lstb.Logger.Println(err)
				return err
			}
			defer stm.Close()
			_, err = stm.Exec(idB, f.Value.(*list.Element).Value.(*users.User).Id)
			if err != nil {
				Lstb.Logger.Println(err)
				return err
			}
			return nil
		})
		if err != nil {
			Lstb.Logger.Println(err)
		}
	}
}

func getIdMessageMax(idBall int64, base *db.Env) (int32, error) {
	var IdMax int32
	rows, err := base.Db.Query("SELECT index_m FROM message WHERE index_m=(SELECT max(index_m) FROM message WHERE containerid=$1);", idBall)
	if err != nil {
		return IdMax, err
	}
	defer rows.Close()
	if rows.Next() != false {
		rows.Scan(&IdMax)
	}
	return IdMax, err
}

/*
	createcontainer(double precision,double precision,double precision,double precision,integer,character varying,integer,date)
	createcontainer(directionc double precision,
	speedc double precision,
	latitudec double precision,
	longitudec double precision,
	idcreatorc integer,
	title character varying,
	idx integer,
	creation date)
*/

func (Lst_ball *All_ball) InsertBallon(eball *list.Element, base *db.Env) (executed bool, err error) {
	NewBall := eball.Value.(*Ball)
	var IdC int64
	err = base.Db.QueryRow("SELECT insertcontainer($1, $2, $3, $4, $5, $6, $7, $8)",
		NewBall.Creator.Value.(*users.User).Id,
		NewBall.Coord.Value.(Checkpoint).Coord.Lat,
		NewBall.Coord.Value.(Checkpoint).Coord.Lon,
		NewBall.Wind.Degress, NewBall.Wind.Speed,
		strings.Trim(NewBall.Title, "\x00"),
		NewBall.Id_ball,
		NewBall.Stats.CreationDate).Scan(&IdC)
	if err != nil {
		Lst_ball.Logger.Println("Error on QueryRow: ", err)
		return false, err
	}
	Lst_ball.SetStatsBallon(IdC, NewBall.Stats, base.Db)
	err = Lst_ball.InsertMessages(NewBall.Messages, IdC, base)
	if err != nil {
		Lst_ball.Logger.Println("Insert Ball fail")
		return false, err
	}
	Lst_ball.SetFollowerBalls(NewBall, base)
	Lst_ball.SetItinerary(base.Db, eball)
	executed = true
	return executed, err
}

func (b *Ball) UpdateLocation(base *db.Env) error {
	var idB int64
	base.Db.QueryRow("SELECT id FROM container WHERE ianix=$1;", b.Id_ball).Scan(&idB)
	id := int64(0)
	if b.Possessed != nil {
		id = b.Possessed.Value.(*users.User).Id
	}

	err := base.Transact(base.Db, func(tx *sql.Tx) error {
		stm, err := tx.Prepare(" SELECT setdatacontainer($1, $2, $3, $4, $5, $6, $7, $8)")
		if err != nil {
			return err
		}
		defer stm.Close()
		_, err = stm.Exec(b.Wind.Degress, b.Wind.Speed, b.Coord.Value.(Checkpoint).Coord.Lat, b.Coord.Value.(Checkpoint).Coord.Lon, idB, id, b.Coord.Value.(Checkpoint).Date, b.Coord.Value.(Checkpoint).MagnetFlag)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

/*
CREATE OR REPLACE FUNCTION public.insertcontainer(idcreatorc integer, latitudec double precision, longitudec double precision, directionc double precision, speedc double precision, title text, idx integer, creation date)
 RETURNS SETOF integer
 LANGUAGE plpgsql
AS $function$  BEGIN RETURN QUERY INSERT INTO container (direction, speed, location_ct, idcreator, titlename, ianix, creationdate) VALUES(directionc, speedc , ST_SetSRID(ST_MakePoint(latitudec, longitudec), 4326), idcreatorc, title, idx, creation) RETURNING id;  END; $function$
\*/

func (ball *Ball) addMessage(base *db.Env) error {
	var idB int64
	base.Db.QueryRow("SELECT id FROM container WHERE ianix=$1;", ball.Id_ball).Scan(&idB)
	idMessageMax, er := getIdMessageMax(idB, base)
	if er != nil {
		return er
	}
	for f := ball.Messages.Front(); f != nil; f = f.Next() {
		mes := f.Value.(Message)
		if mes.Id > idMessageMax {
			err := base.Transact(base.Db, func(tx *sql.Tx) error {
				stm, err := tx.Prepare("INSERT INTO message(content, containerid, index_m, size) VALUES ($1, (SELECT id FROM container WHERE ianix=$2), $3, $4)")
				if err != nil {
					return err
				}
				defer stm.Close()
				_, err = stm.Exec(strings.Trim(mes.Content, "\x00"), ball.Id_ball, f.Value.(Message).Id, f.Value.(Message).Size)
				return err
			})
			if err != nil {
				return er
			}
		}
	}
	return nil
}

func (Lb *All_ball) Update_balls(ABalls *All_ball, base *db.Env) (er error) {
	for e := ABalls.Blist.Front(); e != nil; e = e.Next() {
		ball := e.Value.(*Ball)
		var idB int64
		base.Db.QueryRow("SELECT id FROM container WHERE ianix=$1;", ball.Id_ball).Scan(&idB)
		ball.Lock()
		if ball.FlagC == true {
			Lb.InsertBallon(e, base)
		} else if ball.Edited == true {
			Lb.SetStatsBallon(idB, ball.Stats, base.Db)
			Lb.SetItinerary(base.Db, e)
			er := ball.addMessage(base)
			if er != nil {
				Lb.Logger.Println(er)
			}
			er = e.Value.(*Ball).UpdateLocation(base)
			if er != nil {
				Lb.Logger.Println(er)
			}
			Lb.SetFollowerBalls(e.Value.(*Ball), base)
		}
		ball.Edited = false
		ball.FlagC = false
		ball.Unlock()
	}
	return nil
}
