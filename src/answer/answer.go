// Header

package answer

import (
	"container/list"
	"errors"
	"fmt"
	"protocol"
	"time"
	"users"
)

// History is a history requete that client make to the server
// Log is a last connection requete sock
func Check_user(Req *list.Element, Lst_users *users.All_users) (usr *users.User, err error) {
	user := Lst_users.Lst_users.Front()
	var device *list.Element

	for user != nil {
		device = user.Value.(users.User).Device.Front()
		for device != nil && device.Value.(users.Device).IdMobile != Req.Value.(protocol.Lst_req_sock).IdMobile {
			device = device.Next()
		}
		if device != nil {
			break
		}
		user = user.Next()
	}
	if user == nil {
		usr = new(users.User)
		hist_device := new(users.Device)
		usr.Device = list.New()
		usr.Log = time.Now()
		hist_device.IdMobile = Req.Value.(protocol.Lst_req_sock).IdMobile
		hist_device.History_req = list.New()
		hist_device.History_req.PushFront(users.History{time.Now(), Req.Value.(protocol.Lst_req_sock).Type})
		usr.Device.PushFront(hist_device)
		Lst_users.Lst_users.PushBack(usr)
		Lst_users.Add_new_user(usr)
	} else {
		device.Value.(users.Device).History_req.PushFront(users.History{time.Now(), Req.Value.(protocol.Lst_req_sock).Type})
		usr.Log = time.Now()
	}
	return usr, nil
}

func Check_packets_list(Req *list.Element) bool {
	next := Req.Next()
	tmp := Req
	var nr protocol.Lst_req_sock
	var tr protocol.Lst_req_sock

	for next != nil {
		nr = next.Value.(protocol.Lst_req_sock)
		tr = tmp.Value.(protocol.Lst_req_sock)
		if tr.NbrPack == tr.NumPack-1 {
			return true
		} else if tr.Type != nr.Type {
			return false
		} else if tr.NumPack != nr.NumPack+1 {
			return false
		}
		next = next.Next()
		tmp = tmp.Next()
	}
	if tr.NbrPack == tr.NumPack-1 {
		return true
	}
	return false
}

func Del_request_done(Lst_req *list.List) {
	elem := Lst_req.Front()
	for elem != nil {
		if elem.Value.(protocol.Lst_req_sock).NumPack == elem.Value.(protocol.Lst_req_sock).NbrPack-1 {
			return
		}
		tmp := elem
		Lst_req.Remove(tmp)
		elem = elem.Next()
	}
}

func Manage_type_1(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

func Manage_type_2(Req *list.Element, usr *users.User) (answer []byte) {

	return answer
}

func Manage_type_3(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

func Manage_type_4(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

func Manage_type_5(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

func Get_answer(Lst_req *list.List, Lst_usr *users.All_users) (answer []byte, err error) {
	Req := Lst_req.Front()
	if Req == nil {
		fmt.Println("Error get answer")
	}
	usr, err := Check_user(Req, Lst_usr)
	if err != nil {
		fmt.Println(err)
	} else {
		if Check_packets_list(Req) == true {
			switch Req.Value.(protocol.Lst_req_sock).Type {
			case 1:
				answer = Manage_type_1(Req, usr)
			case 2:
				answer = Manage_type_2(Req, usr)
			case 3:
				answer = Manage_type_3(Req, usr)
			case 4:
				answer = Manage_type_4(Req, usr)
			case 5:
				answer = Manage_type_5(Req, usr)
			}
		}
		Del_request_done(Lst_req)
		return answer, nil
	}
	Del_request_done(Lst_req)
	err = errors.New("Failed check user")
	return answer, err
}

func Get_aknowledgement(Lst_req *list.List, Lst_usr *users.All_users) (answer []byte, err error) {
	return answer, err
}
