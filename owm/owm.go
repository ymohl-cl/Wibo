//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  owm.go                                             :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  By: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  Created: 2015/06/11 13:13:28 by ymohl-cl          #+#    #+#              #
//#  Updated: 2015/06/11 13:13:28 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

/* Package to get wind data from api.openweathermap.org */

package owm

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sync"
)

type Wind struct {
	Speed   float64 `json:"speed"`
	Degress float64 `json:"deg"`
}

type Coordinates struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

/* Station data */
type Weather_data struct {
	Station_id   int         `json:"id"`
	Station_name string      `json:"name"`
	Coord        Coordinates `json:"coord"`
	Wind         Wind        `json:"wind"`
}

/* All data weather map */
type All_data struct {
	Tab_wd []Weather_data `json:"list"`
	sync.RWMutex
	Logger *log.Logger
}

/*
** Source: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
 */
func hsin(theta float64) (result float64) {
	result = math.Pow(math.Sin(theta/2), 2)
	return
}

/*
** Source: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
 */
func (Wd Weather_data) GetDistance(lon_user float64, lat_user float64) float64 {
	var lat1, lat2, lon1, lon2, rayon float64

	lat1 = lat_user * math.Pi / 180
	lon1 = lon_user * math.Pi / 180
	lat2 = Wd.Coord.Latitude * math.Pi / 180
	lon2 = Wd.Coord.Longitude * math.Pi / 180
	rayon = 6378137

	hvsin := hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*hsin(lon2-lon1)
	return (2 * rayon * math.Asin(math.Sqrt(hvsin))) / 1000
}

func (Data *All_data) GetNearest(Lon float64, Lat float64) Weather_data {
	/**
	 ** To BETA TEST
	 **/
	/**
		for _, elem := range Tab_wd.Tab_wd {
			if elem.Station_name == "Paris" {
				station = elem
				break
			}
		}
		return station
	**/
	var es Weather_data
	Data.Lock()
	var best float64
	var save Weather_data

	best = 0.0
	defer Data.Unlock()
	for _, es = range Data.Tab_wd {
		dist := es.GetDistance(Lon, Lat)
		if best == 0.0 || best > dist {
			best = dist
			save = es
		}
	}
	return save
}

/*
** For BetaVersionONLY
 */
/*func (Tab_wd *All_data) Get_Paris() (station Weather_data) {
	for _, elem := range Tab_wd.Tab_wd {
		if elem.Station_name == "Paris" {
			station = elem
			break
		}
	}
	return station
}*/

/* Update_weather_data with api openWeatherMap */
func (Tab_wd *All_data) Update_weather_data() error {
	Tab_wd.Lock()
	defer Tab_wd.Unlock()
	resp, er := http.Get(`http://api.openweathermap.org/data/2.5/box/city?bbox=-90,-180,90,180,10&cluster=yes&APPID=7b7c6c485a78ba0ceaacb887692a88ce`)
	if er != nil {
		return er
	}
	defer resp.Body.Close()
	er = json.NewDecoder(resp.Body).Decode(&Tab_wd)
	if er != nil {
		return er
	}
	return nil
}

/* Print wins data for debug */
func (Wd *All_data) Print_weatherdata() {
	var index int = 0
	for _, elem := range Wd.Tab_wd {
		index++
		Wd.Logger.Println(elem.Station_id)
		Wd.Logger.Println(elem.Station_name)
		Wd.Logger.Println(elem.Coord)
		Wd.Logger.Println(elem.Wind)
		Wd.Logger.Println("--------------------")
	}
	Wd.Logger.Println("Number station updated: %d", index)
}
