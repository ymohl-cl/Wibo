//header

package users

import (
	"container/list"
	//	"fmt"
	"time"
)

/*
** Log is a last signal to device
 */

type History struct {
	Date            time.Time
	Type_req_client int16
}

type Device struct {
	IdMobile    int64
	History_req *list.List //type history
}

type User struct {
	Device *list.List //type Device
	Log    time.Time  //time.Time
}

type All_users struct {
	Lst_users *list.List
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

func (Lst_users *All_users) Del_user(del_user *User) {
	// Supprime un user de la base
	return
}

func (Lst_users *All_users) Add_new_user(new_user *User) {
	// Rajoute un nouvel utilisateur dans la base de donnee.
	return
}

func (Lst_users *All_users) Print_users() {
	// Print All_users
	return
}

func (Lst_users *All_users) Get_users() error {
	// Get all users
	return nil
}
