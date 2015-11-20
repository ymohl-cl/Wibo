package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type Env struct {
	Db *sql.DB
}

/**
* OpenCo
* Open a connection postgres with docker database without ssl.
* returns a Env reference for a data source.
* param: error
* return *Env, error
**/
func (dbp *Env) OpenCo(args []string) (*sql.DB, error) {
	var err error
	str_connect := "user=" + args[3] + " password=" + args[2] + " dbname=" + args[1] + " sslmode=disable host=localhost port=" + args[4]
	db, err := sql.Open("postgres", str_connect)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("No db.driver found")
	}
	dbp.Db = db
	return db, err
}

func (dbp *Env) Transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	return txFunc(tx)
}
