package ballonwork

import (
	"Wibo/db"
	"Wibo/protocol"
	"container/list"
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
	Wlist *list.List
}

func (workBall *WorkBall) Check_neirbycoord(request *list.Element) bool {
	rlon := request.Value.(*protocol.Request).Coord.Lon
	rlat := request.Value.(*protocol.Request).Coord.Lat
	if workBall.Coord != nil {
		coord := workBall.Coord
		if coord.Lon < rlon+0.01 &&
			coord.Lon > rlon-0.01 &&
			coord.Lat < rlat+0.01 &&
			coord.Lat > rlat-0.01 {
			return true
		}
	}
	return false
}

/******************************************************************************/
/******************************** MERGE JAIME *********************************/
/******************************************************************************/

func GetCoord(position string) *Coordinates {
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
	long, err := strconv.ParseFloat(point[0], 15)
	if err != nil {
		return nil
	}
	var lat float64
	lat, _ = strconv.ParseFloat(point[1], 15)
	return &Coordinates{Lon: long, Lat: lat}
}

func (wlist *All_work) Get_workBall(base *db.Env) error {

	wlist.Wlist = list.New()
	rows, err := base.Db.Query("SELECT title, message, ST_AsText(ballonwork.location_wk), link FROM  ballonwork;")
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() != false {
		for rows.Next() {
			var til, ms, pos, url string
			err = rows.Scan(&til, &ms, &pos, &url)
			if err != nil {
				return err
			}
			cord := GetCoord(pos)
			wlist.Wlist.PushBack(&WorkBall{Title: til, Message: ms, Coord: cord, Link: url})
		}
	}
	return nil
}
