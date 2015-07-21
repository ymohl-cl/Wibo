// Header

package answer

import (
	"container/list"
	"errors"
	"fmt"
	"sock"
	"time"
	"users"
)

// History is a history requete that client make to the server
// Log is a last connection requete sock
func Check_user(Req *list.Element, Lst_users *users.All_users) (usr *User, err error) {
	user := Lst_users.Front()
	var device users.Device

	for user != nil {
		device = User.Value.(users.User).Device.Front()
		for device != nil && device.Value.(users.History).IdMobile != Req.Value.(Lst_req_sock).IdMobile {
			device = device.Next()
		}
		if device != nil {
			break
		}
		user = user.Next()
	}
	if user == nil {
		usr = new(users.User)
		hist_device = new(users.Device)
		usr.Device = list.New()
		usr.Log = time.Now()
		hist_device.IdMobile = Req.Value.(Lst_req_sock).IdMobile
		hist_device.history_req = list.New()
		hist_device.history_req.PushFront(users.History{time.Now(), Req.Value.(Lst_req_sock).Type})
		usr.Device.PushFront(hist_device)
		Lst_users.PushBack(usr)
		Lst_users.Add_new_user(usr)
	} else {
		device.history_req.PushFront(users.History{time.Now(), Req.Value.(Lst_req_sock).Type})
		usr = user.Value.(users.User)
		usr.Log = time.Now()
	}
	return usr, nil
}

func Check_packets_list(Req *list.Element) {
	next = Req.Next()
	tmp = Req

	for next != nil {
		nr = next.Value.(Lst_req_sock)
		tr = tmp.Value.(Lst_req_sock)
		if tr.NbrPacket == tr.NumPacket-1 {
			return true
		} else if tr.Type != nr.Type {
			return false
		} else if tr.NumPacket != nr.NumPacket+1 {
			return false
		}
		next = next.Next()
		tmp = tmp.Next()
	}
	if tr.NbrPacket == tr.NumPacket-1 {
		return true
	}
	return false
}

func Del_request_done(Lst_req *list.List) {
	elem = Lst_req.Front()
	for elem != nil {
		if elem.Value.(Lst_req_sock).NumPack == elem.Value.(Lst_req_sock).NbrPack-1 {
			return
		}
		tmp := elem
		Lst_req.Remove(tmp)
		elem = elem.Next()
	}
}

func Manage_type_1(Req *list.Element) (answer []byte) {
	return answer
}

func Manage_type_2(Req *list.Element) (answer []byte) {

	return answer
}

func Manage_type_3(Req *list.Element) (answer []byte) {
	return answer
}

func Manage_type_4(Req *list.Element) (answer []byte) {
	return answer
}

func Manage_type_5(Req *list.Element) (answer []byte) {
	return answer
}

func Get_answer(Lst_req *list.List, Lst_usr *users.All_users) (answer []byte, err error) {
	Req = Lst_req.Front()
	if Req == nil {
		fmt.Println("Error get answer")
	}
	usr, err := Check_user(Req, Lst_usr)
	if err != nil {
		fmt.Println(err)
	} else {
		if Req.Value.(Lst_req_sock).NbrPack > 1 {
			if Check_packets_list(Req) == false {
				break
			}
		}
		switch Req.Value.(Lst_req_sock).Type {
		case 1:
			answer = Manage_type_1(Req)
		case 2:
			answer = Manage_type_2(Req)
		case 3:
			answer = Manage_type_3(Req)
		case 4:
			answer = Manage_type_4(Req)
		case 5:
			answer = Manage_type_5(Req)
		}
		Del_request_done(Lst_req)
		return answer, nil
	}
	Del_request_done(Lst_req)
	err = errors.New("Failed check user")
	return answer, err
}

func Get_aknowledgement(Lst_req *list.List, Lst_usr *users.All_users) (answer []byte) {
	return answer
}
