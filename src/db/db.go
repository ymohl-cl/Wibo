package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
func (dbp *Env) OpenCo(error) (*sql.DB, error) {

	var err error
	db, err := sql.Open("postgres", "user=wibo  password='wibo' dbname=wibo_base sslmode=disable host=localhost port=49155")
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("No db.driver found")
	}
	dbp.Db = db
	return db, err
}

/**
* PingMyBase
* Ping on Database returns without error.
* If ping returns an error that means the connection to the Database is not existing at all.
* Which will require further error checking steps.
* return bool and err message
 */
func (dbp *Env) PingMyBase(Db *sql.DB) (connected bool, err error) {
	if err := Db.Ping(); err != nil {
		return false, err
	}
	return true, err
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

func (dbp *Env) BeginTr() (tx *sql.Tx) {
	tx, err := dbp.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}
