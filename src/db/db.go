package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

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

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

type Handler struct {
	*Env
	hand func(e *Env, w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.hand(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// Retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func GetIndex(env Env, w http.ResponseWriter, r *http.Request) error {
	/*	users, err := env.Db.GetAllUsers()
		if err != nil {
			// We return a status error here, which conveniently wraps the error
			// returned from our DB queries.
			return StatusError{500, err}
		}

		fmt.Fprintf(w, "%+v", users)*/
	return nil
}

/**
* Add OpenCo sql.DB
**/
func (m *Env) OpenCo(error) (*sql.DB, error) {

	// Open a connection.
	var err error
	Db, err := sql.Open("postgres", "user=wibo  password='wibo' dbname=wibo_base sslmode=disable host=localhost port=49155")
	if err != nil {
		return Db, errors.New("No db.driver found")
	}
	defer Db.Close()
	m.Db = Db
	return Db, err
}
