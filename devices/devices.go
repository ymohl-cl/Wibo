package devices

import (
	"Wibo/protocol"
	"Wibo/users"
	"container/list"
	"database/sql"
	"errors"
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

func (Devices *All_Devices) GetDevice(request *list.Element, Db *sql.DB, Data *answer.Data) (dvc *list.Element, er error) {
	req := request.Value.(protocol.Request)
	ed := Devices.Dlist.Front()
	er = nil

	if req.IdMobile.Len() == 1 {
		er = errors.New("Id mobile bad format")
		return nil, er
	}
	for ed != nil && Compare(ed.Value.(*Device).Id, req.IdMobile) != 0 {
		ed = ed.Next()
	}
	if ed == nil {
		ed = Devices.AddDeviceOnBdd(req.IdMobile, Data.Lst_users, Db)
	}
	return ed, er
}

func (Device *Device) AddUserSpecOnHistory(euser *list.Element) {
	e := Device.Historic
	user1 := euser.Value.(*users.User)

	for e != nil {
		user2 := e.Value.(*list.Element).(*users.User)
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
** Cette fonction cree un Device et un user par default.
** Ajouter l'user par default dans la liste des users.
** Ajouter le list element de l'user dans Device.UserDefault.
 */
func (Devices *All_Devices) AddDeviceOnBdd(Id string, Ulist *users.All_users, Db *sql.DB) *list.Element {
	newDevice := new(Device)
	newDevice.Historic = list.New()
	newDevice.Id = Id
	newDevice.UserDefault = Ulist.AddNewDefaultUser()
	// AddNewDefaultUser: creer un utilisateur par default pour le device.
	// L'ajoute a la liste des users.
	// Insere l'user dans la base de donnee.
	//userspec a null
	//Historic  = list new

	// ICI AJOUTER LE DEVICE A LA BASE DE DONNEE DES DEVICES.
	Ed := Devices.Dlist.PushFront(newDevice)
	return Ed
}
