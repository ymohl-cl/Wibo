package ballonwork

import (
	"Wibo/db"
	"Wibo/protocol"
	"container/list"
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
func (wlist *All_work) Get_workBall(base *db.Env) error {
	return nil
}
