//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  users.go                                           :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  by: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

package users

import (
	"container/list"
	"protocol"
	"time"
)

/* History request */
type History struct {
	Date            time.Time
	Type_req_client int16
}

type Device struct {
	IdMobile    int64      /* type int64 is temporary */
	History_req *list.List /* Value: History */
}

type User struct {
	Id       int64
	Device   *list.List /* Value: Device */
	Log      time.Time  /*Date of the last query */
	Followed *list.List /* Value: *list.Element.Value.(*ballon.Ball) */
}

type All_users struct {
	Ulist  *list.List
	Id_max int64
}

func (User *User) User_is_online() bool {
	t_now := time.Now()
	t_user := User.Log
	if t_user.Hour() == t_now.Hour() && t_user.Minute() > t_now.Minute()-2 {
		return true
	} else {
		return false
	}
}

/*
** Manage users's connexion
 */
func (ulist *All_users) Check_user(request *list.Element) (user *list.Element, err error) {
	user = ulist.Ulist.Front()
	var device *list.Element

	rqt := request.Value.(protocol.Request)
	for user != nil {
		device = user.Value.(*User).Device.Front()
		for device != nil && device.Value.(Device).IdMobile != rqt.Deviceid {
			device = device.Next()
		}
		if device != nil {
			break
		}
		user = user.Next()
	}
	if user == nil {
		usr := new(User)
		var hist_device Device
		usr.Device = list.New()
		usr.Log = time.Now()
		usr.Followed = list.New()
		hist_device.IdMobile = request.Value.(protocol.Request).Deviceid
		hist_device.History_req = list.New()
		hist_device.History_req.PushFront(History{time.Now(), request.Value.(protocol.Request).Rtype})
		usr.Device.PushFront(hist_device)
		user = ulist.Ulist.PushBack(usr)
		ulist.Add_new_user(usr)
	} else {
		device.Value.(Device).History_req.PushFront(History{time.Now(), request.Value.(protocol.Request).Rtype})
		user.Value.(*User).Log = time.Now()
	}
	return user, nil
}

/* Delete user from database */ /* Why ? xd */
func (Lst_users *All_users) Del_user(del_user *User) {
	return
}

/* Add user in database */
func (Lst_users *All_users) Add_new_user(new_user *User) {
	return
}

/* Print_users for debug */
func (Lst_users *All_users) Print_users() {
	return
}

/*
** Get_users va recuperer tous les utilisateurs dans la base de donnee.
 */
func (Lst_users *All_users) Get_users() error {
	return nil
}
