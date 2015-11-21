package ballonwork

import (
	"Wibo/db"
	"Wibo/protocol"
	"container/list"
	"log"
	"math"
	"strconv"
	"strings"
)

type Coordinates struct {
	Lon float64
	Lat float64
}

type WorkBall struct {
	Title   string       //16 octets
	Message string       // 664 octets
	Coord   *Coordinates // 16 octets
	Link    string       //320 octets
}

type All_work struct {
	Wlist  *list.List
	Logger *log.Logger
}

func (ball *WorkBall) Check_nearbycoord(request *list.Element) bool {
	rlon := request.Value.(*protocol.Request).Coord.Lon
	rlat := request.Value.(*protocol.Request).Coord.Lat

	if ball.GetDistance(rlon, rlat) < 1.0 {
		return true
	}
	return false
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
func (ball *WorkBall) GetDistance(lon_user float64, lat_user float64) float64 {
	var lat1, lat2, lon1, lon2, rayon float64
	lat1 = lat_user * math.Pi / 180
	lon1 = lon_user * math.Pi / 180
	lat2 = ball.Coord.Lat
	lon2 = ball.Coord.Lon
	lat2 = ball.Coord.Lat * math.Pi / 180
	lon2 = ball.Coord.Lon * math.Pi / 180
	rayon = 6378137.0

	hvsin := hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*hsin(lon2-lon1)
	return (2 * rayon * math.Asin(math.Sqrt(hvsin))) / 1000
}

/******************************************************************************/
/******************************** MERGE JAIME *********************************/
/******************************************************************************/

func (Wlist *All_work) GetCoord(position string) *Coordinates {
	// Return true if 'value' char.
	f := func(c rune) bool {
		return c == '(' || c == '(' || c == ')' || c == '"' ||
			c == 'P' || c == 'O' || c == 'I' || c == 'N' ||
			c == 'T'
	}
	// Separate into fields with func.
	fields := strings.FieldsFunc(position, f)
	// Separate into cordinates  with Fields.
	point := strings.Fields(fields[0])
	lat, err := strconv.ParseFloat(point[0], 15)
	if err != nil {
		Wlist.Logger.Println("Error ParseFloat on GetCoord: ", err)
		return nil
	}
	long, err := strconv.ParseFloat(point[1], 15)
	if err != nil {
		Wlist.Logger.Println("Error ParseFloat on GetCoord: ", err)
		return nil
	}
	return &Coordinates{Lon: long, Lat: lat}
}

func (wlist *All_work) Get_workBall(base *db.Env) error {
	rows, err := base.Db.Query("SELECT title, message, ST_AsText(ballonwork.location_wk), link FROM  ballonwork;")
	if err != nil {
		wlist.Logger.Println("Error Query on Get_workBall: ", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var til, ms, pos, url string
		err = rows.Scan(&til, &ms, &pos, &url)
		if err != nil {
			wlist.Logger.Println("Error Scan on Get_workBall: ", err)
			return err
		}
		cord := wlist.GetCoord(pos)
		wlist.Wlist.PushBack(&WorkBall{Title: til, Message: ms, Coord: cord, Link: url})
	}
	wlist.Logger.Println("Get number WorkBall: ", wlist.Wlist.Len())
	return nil
}
