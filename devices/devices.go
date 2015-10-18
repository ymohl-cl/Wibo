package devices

import (
	"Wibo/protocol"
	"Wibo/users"
	"container/list"
	"database/sql"
	"errors"
	"strings"
)

/* *list.Element.Value.(*users.User) */
/* Historic to user historic logged on this device */
type Device struct {
	Idbdd         int64
	Id            string
	IdUserDefault int64
	UserDefault   *list.Element
	UserSpec      *list.Element
	Historic      *list.List
}

type All_Devices struct {
	Dlist *list.List
}

func (Devices *All_Devices) GetDevice(request *list.Element, Db *sql.DB, Ulist *users.All_users) (dvc *list.Element, er error) {
	req := request.Value.(protocol.Request)
	ed := Devices.Dlist.Front()
	er = nil

	if len(req.IdMobile) == 1 {
		er = errors.New("Id mobile bad format")
		return nil, er
	}
	for ed != nil && strings.Compare(ed.Value.(*Device).Id, req.IdMobile) != 0 {
		ed = ed.Next()
	}
	if ed == nil {
		ed, er = Devices.AddDeviceOnBdd(req.IdMobile, Ulist, Db)
	}
	return ed, er
}

func (Device *Device) AddUserSpecOnHistory(euser *list.Element) {
	e := Device.Historic.Front()
	user1 := euser.Value.(*users.User)

	for e != nil {
		user2 := e.Value.(*list.Element).Value.(*users.User)
		if user2.Id == user1.Id {
			break
		}
		e = e.Next()
	}
	if e == nil {
		Device.Historic.PushFront(euser)
	}
}

/******************************************************************************/
/********************************* MERGE JAIME ********************************/
/******************************************************************************/

func (Devices *All_Devices) AddDeviceOnBdd(Id string, Ulist *users.All_users, Db *sql.DB) (*list.Element, error) {
	var err error
	newDevice := new(Device)
	newDevice.Historic = list.New()
	newDevice.Id = Id
	newDevice.UserDefault = Ulist.AddNewDefaultUser(Db)
	if err != nil {
		return nil, err
	}
	newDevice.UserSpec = nil
	rows, err := Db.Query("INSERT INTO device (id_type_d, typename, idclient, user_id_user) VALUES ($1, $2, $3, $4) RETURNING id;", 1, "device_default", newDevice.Id, newDevice.UserDefault.Value.(*users.User).Id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&newDevice.IdUserDefault)
		if err != nil {
			return nil, err
		}
	}
	return Devices.Dlist.PushFront(newDevice), nil
}
