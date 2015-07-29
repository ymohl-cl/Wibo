package main

import (
	"container/list"
	"fmt"
	"github.com/Wibo/src/ballon"
	"github.com/Wibo/src/db"
	"github.com/Wibo/src/users"
	"time"
)

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
	myDb := new(db.Env)
	Db, err := myDb.OpenCo(err)
	checkErr(err)
	// Initialise our app-wide environment with the services/info we need.
	//users.Get_users()
	/*	for {
		fmt.Println("manage server")
		time.Sleep(time.Second * 30)
	}*/
	my_users := new(users.All_users)
	my_users.Lst_users = list.New()
	my_users.Get_users(Db)
	ball := new(ballon.All_ball)
	ball.Lst = list.New()
	newUsr := my_users.SelectUser(2, Db)
	ball.GetListBallsByUser(*newUsr, Db)
	ball.Get_balls(my_users, Db)
	ball.Print_all_balls()
	//my_users.Del_user(newUsr, Db)
}
