//header

package users

import (
	"container/list"
)

/*
** Log is a last signal to device
 */
type User struct {
	Device string
	Log    int
}

type All_users struct {
	Lst_users list.List
}

func (User *User) User_is_online() bool {
	if User.Log == 0 {
		return true
	} else {
		return false
	}
}

func (Lst_users *All_users) Del_user(del_user *User) {
	return
}

func (Lst_users *All_users) Add_new_user(new_user *User) {
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
