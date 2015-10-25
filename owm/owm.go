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
	"fmt"
	"net/http"
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
}

func (Tab_wd *All_data) Get_Paris() (station Weather_data) {
	for _, elem := range Tab_wd.Tab_wd {
		if elem.Station_name == "Paris" {
			station = elem
			break
		}
	}
	return station
}

/* Update_weather_data with api openWeatherMap */
func (Tab_wd *All_data) Update_weather_data() error {
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
func (Tab_wd *All_data) Print_weatherdata() {
	var index int = 0
	for _, elem := range Tab_wd.Tab_wd {
		index++
		fmt.Println(elem.Station_id)
		fmt.Println(elem.Station_name)
		fmt.Println(elem.Coord)
		fmt.Println(elem.Wind)
		fmt.Println("--------------------")
	}
	fmt.Println("nombre de stations: %d", index)
}
