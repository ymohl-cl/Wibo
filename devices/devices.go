package devices

import (
	"Wibo/db"
	"Wibo/protocol"
	"Wibo/users"
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

/* *list.Element.Value.(*users.User) */
/* Historic to user historic logged on this device */
type Device struct {
	Idbdd         int64         /* Idbdd */
	Id            string        /* Chaine de 40 octets pour l'identification unique device */
	IdUserDefault int64         /* Id de la bdd user, du user default */
	UserDefault   *list.Element /* User default device */
	UserSpec      *list.Element /** User specifique compte existant */
	Historic      *list.List    /* Liste de *list.Element.Value.(users.User) qui se sont deja connecte depuis ce device. Le user Default est exclu */
}

type All_Devices struct {
	Dlist *list.List
}

func (Devices *All_Devices) GetDevice(request *list.Element, Db *sql.DB, Ulist *users.All_users) (dvc *list.Element, er error) {
	req := request.Value.(*protocol.Request)
	ed := Devices.Dlist.Front()
	er = nil

	if len(req.IdMobile) == 1 {
		er = errors.New("Id mobile bad format")
		fmt.Println(er)
		return nil, er
	}
	for ed != nil && strings.Compare(ed.Value.(*Device).Id, req.IdMobile) != 0 {
		ed = ed.Next()
	}
	if ed == nil {
		ed, er = Devices.AddDeviceOnBdd(req, Ulist, Db)
		if er != nil {
			fmt.Println("Echec add device on bdd")
		}
	}
	fmt.Println(er)
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

/* 
Idbdd         int64         
	Id            string        
	IdUserDefault int64         
	UserDefault   *list.Element 
	UserSpec      *list.Element 
	Historic      *list.List   
*/
func (dlist *All_Devices) Get_devices(LstU *users.All_users, base *db.Env) error {
	for u := LstU.Ulist.Front(); u != nil; u = u.Next() {
		rows, err := Db.Query("SELECT id, user_id_user FROM device WHERE user_id_user=$1", u.Id)
		if err != nil {
			fmt.Println(err)
		}
		if rows.Next() != false {
			for rows.Next() {
				var tdate time.Time
				var attm bool
				var point string
				rows.Scan(&tdate, &attm, &point)
				fmt.Println(tdate, attm, point)
				fmt.Printf("%T | %v ", point, point)
				Ball.Itenerary.PushBack(&Checkpoint{Date: tdate})
			}
			if err != nil{
				fmt.Println(err)
			}
		}
	}

	return nil
}

func (Devices *All_Devices) AddDeviceOnBdd(req *protocol.Request, Ulist *users.All_users, Db *sql.DB) (*list.Element, error) {
	var err error
	newDevice := new(Device)
	newDevice.Historic = list.New()
	newDevice.Id = req.IdMobile
	newDevice.UserDefault = Ulist.AddNewDefaultUser(Db, req)
	if newDevice.UserDefault == nil {
		return nil, errors.New("Add new default user not permission")
	}
	newDevice.IdUserDefault = newDevice.UserDefault.Value.(*users.User).Id
	if err != nil {
		return nil, err
	}
	newDevice.UserSpec = nil
	rows, err := Db.Query("INSERT INTO device (id_type_d, typename, idclient, user_id_user) VALUES ($1, $2, $3, $4) RETURNING id;", 1, "device_default", newDevice.Id, newDevice.IdUserDefault)
	if err != nil {
		fmt.Println("Db query pas content")
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&newDevice.Idbdd)
		if err != nil {
			fmt.Println("Rows scan pas content")
			fmt.Println(err)
			return nil, err
		}
	}
	return Devices.Dlist.PushFront(newDevice), nil
}
