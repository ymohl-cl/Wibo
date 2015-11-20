package ballon

import (
	"Wibo/db"
	"Wibo/users"
	"container/list"
	"database/sql"
	"fmt"
	af "github.com/spf13/afero"
	"io"
	"io/ioutil"
	"path"
	_ "path/filepath"
	"strings"
)

func removeDir() error {
	for _, fs := range Fss {
		err := fs.RemoveAll(SubDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func newFile(Name string, fs af.Fs) (f af.File) {
	fs.MkdirAll(SubDir, 0777)
	f, err := fs.Create(path.Join(SubDir, Name))
	if err != nil {
		return nil
	}
	return f

}

func writeFile(fs af.Fs, fname string, flag int, text string) string {
	f, err := fs.OpenFile(fname, flag, 0666)
	if err != nil {
	}
	n, err := io.WriteString(f, text)
	fmt.Println(n)
	if err != nil {

	}
	f.Close()
	data, err := ioutil.ReadFile(fname)
	if err != nil {

	}
	return string(data)

}

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
			row, err := stm.Query(e.Value.(Message).Content, idBall, i, e.Value.(Message).Size)
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
	fmt.Printf("id %v\n", idB)
	base.Db.QueryRow("DELETE FROM followed WHERE container_id=$1", idB)
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
	err := base.Transact(base.Db, func(tx *sql.Tx) error {
		var err error
		stm, err := tx.Prepare("select index_m from message where index_m = (select max(index_m) from message) and containerid=$1;")
		if err != nil {
			return err
		}
		defer stm.Close()
		rs, err := stm.Query(idBall)
		if err != nil {
			return err
		}
		defer rs.Close()
		if rs.Next() != false {
			rs.Scan(&IdMax)
		}
		return err
	})
	fmt.Println("IdMax message: ", IdMax)
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

func (Lst_ball *All_ball) InsertBallon(NewBall *Ball, base *db.Env) (executed bool, err error) {
	var IdC int64
	err = base.Db.QueryRow("SELECT insertcontainer($1, $2 , $3, $4, $5, $6, $7, $8)",
		NewBall.Creator.Value.(*users.User).Id,
		NewBall.Coord.Value.(Checkpoint).Coord.Lat,
		NewBall.Coord.Value.(Checkpoint).Coord.Lon,
		NewBall.Wind.Degress, NewBall.Wind.Speed,
		strings.Trim(NewBall.Title, "\x00"),
		NewBall.Id_ball,
		NewBall.Stats.CreationDate).Scan(&IdC)
	//Lst_ball.Logger.Println(err)
	if err != nil {
		return false, err
	}
	Lst_ball.checkErr(err)
	Lst_ball.SetStatsBallon(IdC, NewBall.Stats, base.Db)
	NewBall.SetCreationCoordOnItinerary(base.Db, Lst_ball.Logger)
	err = Lst_ball.InsertMessages(NewBall.Messages, IdC, base)
	if err != nil {
		Lst_ball.Logger.Println("Insert Ball fail")
		return false, err
	}
	Lst_ball.SetFollowerBalls(NewBall, base)
	executed = true
	return executed, err
}

func (b *Ball) UpdateLocation(base *db.Env) error {
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
		_, err = stm.Exec(b.Wind.Degress, b.Wind.Speed, b.Coord.Value.(Checkpoint).Coord.Lat, b.Coord.Value.(Checkpoint).Coord.Lon, b.Id_ball, id, b.Coord.Value.(Checkpoint).Date, b.Coord.Value.(Checkpoint).MagnetFlag)
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
	idMessageMax, er := getIdMessageMax(ball.Id_ball, base)
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
				_, err = stm.Exec(mes.Content, ball.Id_ball, f.Value.(Message).Id, f.Value.(Message).Size)
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
		ball.Lock()
		if ball.FlagC == true {
			Lb.InsertBallon(e.Value.(*Ball), base)
		} else if ball.Edited == true {
			Lb.SetStatsBallon(ball.Id_ball, ball.Stats, base.Db)
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

/*func (Lb *All_ball) Update_balls(ABalls *All_ball, base *db.Env) (er error) {
	fmt.Println("Enter update")
	for e := ABalls.Blist.Front(); e != nil; e = e.Next() {
		fmt.Println("title ball: ", e.Value.(*Ball).Title)
		fmt.Println("FlagcC ", e.Value.(*Ball).FlagC)
		fmt.Println("Edited ", e.Value.(*Ball).Edited)
		if e.Value.(*Ball).FlagC == true {
			fmt.Println("try insert: ")
			Lb.InsertBallon(e.Value.(*Ball), base)
			e.Value.(*Ball).Edited = false
			e.Value.(*Ball).FlagC = false
		} else if e.Value.(*Ball).Edited == true {
			fmt.Println("try update: ")
			e.Value.(*Ball).Lock()
			idBall := e.Value.(*Ball).Id_ball
			idMessageMax, er := getIdMessageMax(idBall, base)
			if er != nil {
				Lb.Logger.Println(er)
			}
			Lb.SetStatsBallon(idBall, e.Value.(*Ball).Stats, base.Db)
			fmt.Println("Set itinerary ball id: ", idBall)
			Lb.SetItinerary(base.Db, e)
			for f := e.Value.(*Ball).Messages.Front(); f != nil; f = f.Next() {
				if f.Value.(Message).Id > idMessageMax {
					err := base.Transact(base.Db, func(tx *sql.Tx) error {
						stm, err := tx.Prepare("INSERT INTO message(content, containerid, index_m, size) VALUES ($1, (SELECT id FROM container WHERE ianix=$2), $3, $4)")
						if err != nil {
							return err
						}
						defer stm.Close()
						_, err = stm.Exec(f.Value.(Message).Content, idBall, f.Value.(Message).Id, f.Value.(Message).Size)
						if err != nil {
							return err
						}
						return err
					})
					if err != nil {
						Lb.Logger.Println(err)
					}
				}
			}
			er = e.Value.(*Ball).UpdateLocation(base)
			if er != nil {
				Lb.Logger.Println(er)
			}
			Lb.SetFollowerBalls(e.Value.(*Ball), base)
			e.Value.(*Ball).Unlock()
		}
	}
	return er
}*/
