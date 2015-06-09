package main

import "fmt"
import "time"

func	calcul_ballon() {
	fmt.Println("Calcul position ball")
}

func	check_weather_data() {
	fmt.Println("Get weather data")
}

func	test1() {
	for {
		time.Sleep(time.Second * 100)
		calcul_ballon()
	}
}

func	test2() {
	for {
		time.Sleep(time.Second * 100)
		check_weather_data()
	}
}

func main() {
	go test1()
	go test2()
	for {
		fmt.Println("manage server")
		time.Sleep(time.Second * 30)
	}
}
