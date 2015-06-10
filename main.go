package main

import (
	"encoding/json"
	"fmt"
	//	"io/ioutil"
	"net/http"
	"time"
)

type Coordinates struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type Sys struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

// Wind struct contains the speed and degree of the wind.
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

// Weather struct holds high-level, basic info on the returned
// data.
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type CurrentWeatherData struct {
	GeoPos  Coordinates `json:"coord"`
	Sys     Sys         `json:"sys"`
	Base    string      `json:"base"`
	Weather []Weather   `json:"weather"`
	//	Main    Main        `json:"main"`
	//	Wind    Wind        `json:"wind"`
	//	Clouds  Clouds      `json:"clouds"`
	Dt   int    `json:"dt"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
	Unit string
	Lang string
}

type AllData struct {
	List []CurrentWeatherData `json:"list"`
}

func calcul_ballon() {
	fmt.Println("Calcul position ball")
}

func check_weather_data() {
	fmt.Println("requete\n")
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/box/city?bbox=-180,-90,180,90,10&cluster=yes")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println("fin de la requete\n")
	var list AllData
	//	var toto []CurrentWeatherData
	fmt.Println("Decodage json\n")
	json.NewDecoder(resp.Body).Decode(&list)

	fmt.Println("Print de la list data\n")
	var index int = 0
	for _, w := range list.List {
		index++
		fmt.Println("elem\n")
		//		fmt.Println(w)
		fmt.Printf("%#v\n", w)
		//		fmt.Println(w.List.Base)
		//		fmt.Println(w.List.ID)
		//		fmt.Println(w.List.Name)
		//		fmt.Println(w.List.GeoPos.Longitude)
		//		fmt.Println(w.List.GeoPos.Latitude)
		fmt.Println("--------------------")
		//		fmt.Println("\n")
	}
	fmt.Println("nombre de stations: %d\n", index)
	//	body, err := ioutil.ReadAll(resp.Body)
	//	for _, element := range body {
	//		fmt.Println(element)
	//		fmt.Println("end data")
	//	}

	//	json.NewDecoder(response.Body).Decode(&w);
	//	fmt.Println(resp)
	//	fmt.Println("end data")
	//	fmt.Println(err)
}

func test1() {
	for {
		time.Sleep(time.Second * 100)
		calcul_ballon()
	}
}

func test2() {
	for {
		check_weather_data()
		time.Sleep(time.Second * 100)
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
