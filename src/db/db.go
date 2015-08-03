package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

type Env struct {
	Db *sql.DB
}

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
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

func (dbp *Env) BeginTr() (tx *sql.Tx) {
	tx, err := dbp.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}
