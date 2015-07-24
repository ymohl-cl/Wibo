//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  answer.go                                          :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  By: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  Created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  Updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

package answer

/*
** Ce package est distine a creer la reponse au client en respectant le
** protocole wibo definit dans le trello, rubrique tools.
 */

import (
	"ballon"
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"protocol"
	"time"
	"users"
)

/*
** Check_user verifie l'existance de l'utilisateur dans la base a l'aide de
** son ID_mobile. Si l'utilisateur existe pas, il le rajoute et demande une
** insertion dans la base de donnee.
** Cette fonction retourne l'utilisateur ou error si une anomalie c'est produite
** Check_user ajoute egalement la requete a l'historique des requetes du Device
 */
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

/*
** Check_packets_list Verifie la validite des packets suivants en partant du
** principe que le header du premier packet est valide (traitement multi-packet).
 */
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

/*
** Supprime une requete traite. Gere les multi-paquets
 */
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

func Write_header(Packet []byte, Len int16, Type int16, NbrPack int32, NumPack int32) (answer []byte) {
	Buffer := bytes.NewBuffer(Packet)
	binary.Write(Buffer, binary.BigEndian, Len)
	binary.Write(Buffer, binary.BigEndian, Type)
	Buffer.Next(4)
	binary.Write(Buffer, binary.BigEndian, NbrPack)
	binary.Write(Buffer, binary.BigEndian, NumPack)
	answer = Buffer.Bytes()
	return answer
}

func Write_in_buffer_type1(Req *list.Element, list_tmp *list.List) (answer []byte, Len int16) {
	answer = make([]byte, 1024)
	Buffer := bytes.NewBuffer(answer)

	Buffer.Next(16)
	binary.Write(Buffer, binary.BigEndian, list_tmp.Len())
	Buffer.Next(4)
	elem := list_tmp.Front()
	for elem != nil {
		ball := elem.Value.(ballon.Ball)
		binary.Write(Buffer, binary.BigEndian, ball.Id_ball)
		binary.Write(Buffer, binary.BigEndian, ball.Name)
		Buffer.Next(16 - len(ball.Name))
		binary.Write(Buffer, binary.BigEndian, ball.Coord.Value.(ballon.Checkpoints).Coord.Longitude)
		binary.Write(Buffer, binary.BigEndian, ball.Coord.Value.(ballon.Checkpoints).Coord.Latitude)
		binary.Write(Buffer, binary.BigEndian, ball.Wind.Speed)
		binary.Write(Buffer, binary.BigEndian, ball.Wind.Degress)
	}
	Len = (int16)(list_tmp.Len() * 64) // 64 octet par ballon
	Len += 16                          // 16 octet du header
	answer = Buffer.Bytes()
	return answer, Len
}

/* Manage_type_1 Remplie le buffer avec une reponse au type 1 de la requete avec un maximum de 10 ballons. */
func Manage_type_1(Req *list.Element, usr *users.User, Lst_ball *ballon.All_ball) (answer []byte) {
	list_tmp := list.New()
	elem := Lst_ball.Lst.Front()
	for elem != nil {
		Coord := elem.Value.(ballon.Ball).Coord.Value.(ballon.Checkpoints).Coord
		if Coord.Longitude < Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Longitude+0.01 && Coord.Longitude > Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Longitude-0.01 && Coord.Latitude < Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Latitude+0.01 && Coord.Latitude > Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Latitude-0.01 {
			list_tmp.PushFront(elem.Value.(ballon.Ball))
		}
		elem = elem.Next()
	}
	Len := Lst_ball.Lst.Len()
	for Len > 10 {
		elem := list_tmp.Front()
		random := rand.Intn(Len)
		for elem != nil && random > 0 {
			elem = elem.Next()
			random -= 1
		}
		list_tmp.Remove(elem)
		Len -= 1
	}
	answer, LenPack := Write_in_buffer_type1(Req, list_tmp)
	answer = Write_header(answer, LenPack, 1, 1, 0)
	return answer
}

/* Manage_type_2 Remplie le buffer avec une reponse au type 2 de la requete. */
func Manage_type_2(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

/* Manage_type_3 Remplie le buffer avec une reponse au type 3 de la requete. */
func Manage_type_3(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

/* Manage_type_4 Remplie le buffer avec une reponse au type 4 de la requete. */
func Manage_type_4(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

/* Manage_type_5 Remplie le buffer avec une reponse au type 5 de la requete. */
func Manage_type_5(Req *list.Element, usr *users.User) (answer []byte) {
	return answer
}

/*
** Get_answer fournis une reponse approprie a la requete du client,
** avec un buffer de 1024 Octets. Elle initialisera l'authentification de
** de l'utilisateur et nettoiera le flux de requetes traites.
 */
func Get_answer(Lst_req *list.List, Lst_usr *users.All_users, Lst_ball *ballon.All_ball) (answer []byte, err error) {
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
				answer = Manage_type_1(Req, usr, Lst_ball)
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

/*
** Get_acknowledgement creer un buffer pour confirmation au client et le rempli.
 */
func Get_aknowledgement(Lst_req *list.List, Lst_usr *users.All_users) (answer []byte, err error) {
	return answer, err
}
