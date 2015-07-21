package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"github.com/Wibo/src/db"
	"github.com/Wibo/src/usr"
	//	"net/http"
	"os"
	"time"
)

type Env struct {
	Db   *sql.DB
	Port string
	Host string
}

func calcul_ballon() {
	fmt.Println("Calcul position ball")
}

func check_weather_data() {
	fmt.Println("Get weather data")
}

func test1() {
	for {
		time.Sleep(time.Second * 100)
		calcul_ballon()
	}
}

func test2() {
	for {
		time.Sleep(time.Second * 100)
		check_weather_data()
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//	go test1()
	//	go test2()
	var err error
	Dbl, err := db.OpenCo(err)
	checkErr(err)
	// Initialise our app-wide environment with the services/info we need.
	env := Env{
		Db:   Dbl,
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
	}
	defer Dbl.Close()
	//usr.Get_users()
	/*	for {
		fmt.Println("manage server")
		time.Sleep(time.Second * 30)
	}*/
	//http.Handler("/", db.Handler{env, db.GetIndex})
	lUser := list.New(usr.All_users)
	err = usr.Get_users(env)
}
