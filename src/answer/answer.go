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

type InfoBall struct {
	Id    int64
	Title string
	Lon   float64
	Lat   float64
	Wins  float64
	Wind  float64
}

type Pack6 struct {
	NbrBall int32
	Ifball  *list.List
}

type Message struct {
	Idmess    int32
	Size      int32
	Idcountry int32
	Idcity    int32
	Message   string
	Type_     int32
}

type Pack7 struct {
	Nbruser int32
	NbrMess int32
	Mess    *list.List
}

type Acknow struct {
	Typeack  int32
	Status   int32
	Idball   int64
	IdMobile int64
}

type Header struct {
	NbrOctet   int32
	TypeReq    int32
	NbrPackReq int32
	NumPackReq int32
}

type Packet struct {
	Head  Header
	TPack interface{}
}

type Data struct {
	Lst_req   *list.List
	Lst_asw   *list.List
	Lst_ball  *ballon.All_ball
	Lst_users *users.All_users
	User      users.User
}

/*
** Check_user verifie l'existance de l'utilisateur dans la base a l'aide de
** son ID_mobile. Si l'utilisateur existe pas, il le rajoute et demande une
** insertion dans la base de donnee.
** Cette fonction retourne l'utilisateur ou error si une anomalie c'est produite
** Check_user ajoute egalement la requete a l'historique des requetes du Device
 */
func Check_user(Req *list.Element, Lst_users *users.All_users) (usr users.User, err error) {
	fmt.Println("Check_user")
	user := Lst_users.Lst_users.Front()
	fmt.Println(user)
	var device *list.Element

	for user != nil {
		fmt.Println("user exist")
		device = user.Value.(users.User).Device.Front()
		fmt.Println("device.current")
		for device != nil && device.Value.(users.Device).IdMobile != Req.Value.(protocol.Lst_req_sock).IdMobile {
			fmt.Println("device.Next")
			device = device.Next()
		}
		if device != nil {
			fmt.Println("Break")
			break
		}
		user = user.Next()
		fmt.Println("user.Next")
	}
	fmt.Println("user no exist")
	if user == nil {
		var usr users.User
		var hist_device users.Device
		usr.Device = list.New()
		usr.Log = time.Now()
		hist_device.IdMobile = Req.Value.(protocol.Lst_req_sock).IdMobile
		hist_device.History_req = list.New()
		hist_device.History_req.PushFront(users.History{time.Now(), Req.Value.(protocol.Lst_req_sock).Type})
		usr.Device.PushFront(hist_device)
		Lst_users.Lst_users.PushBack(usr)
		Lst_users.Add_new_user(&usr)
	} else {
		device.Value.(users.Device).History_req.PushFront(users.History{time.Now(), Req.Value.(protocol.Lst_req_sock).Type})
		usr.Log = time.Now()
	}
	fmt.Println("End Check_user")
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

	tr = tmp.Value.(protocol.Lst_req_sock)
	for next != nil {
		nr = next.Value.(protocol.Lst_req_sock)
		if tr.NbrPack == tr.NumPack+1 {
			return true
		} else if tr.Type != nr.Type {
			return false
		} else if tr.NumPack != nr.NumPack+1 {
			return false
		}
		next = next.Next()
		tmp = tmp.Next()
		tr = tmp.Value.(protocol.Lst_req_sock)
	}
	if tr.NbrPack == tr.NumPack+1 {
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

func Write_header(answer Packet) (Buffer *bytes.Buffer) {
	Buffer = new(bytes.Buffer)

	binary.Write(Buffer, binary.BigEndian, answer.Head.NbrOctet)
	binary.Write(Buffer, binary.BigEndian, answer.Head.TypeReq)
	binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
	binary.Write(Buffer, binary.BigEndian, answer.Head.NbrPackReq)
	binary.Write(Buffer, binary.BigEndian, answer.Head.NumPackReq)
	return Buffer
}

func Write_type1(Req *list.Element, list_tmp *list.List) (buf []byte) {
	fmt.Println("Write_type_1")
	var answer Packet
	var typesp Pack6

	answer.Head.NbrOctet = 24
	answer.Head.TypeReq = 6
	answer.Head.NbrPackReq = 1
	answer.Head.NumPackReq = 0
	typesp.NbrBall = (int32)(list_tmp.Len())
	typesp.Ifball = list_tmp
	answer.TPack = typesp
	elem := list_tmp.Front()
	for elem != nil {
		answer.Head.NbrOctet += 56
		elem = elem.Next()
	}
	Buffer := Write_header(answer)
	binary.Write(Buffer, binary.BigEndian, answer.TPack.(Pack6).NbrBall)
	binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
	elem = answer.TPack.(Pack6).Ifball.Front()
	for elem != nil {
		ifb := elem.Value.(InfoBall)
		binary.Write(Buffer, binary.BigEndian, ifb.Id)
		Buffer.WriteString(ifb.Title)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 16-len(ifb.Title)))
		binary.Write(Buffer, binary.BigEndian, ifb.Lon)
		binary.Write(Buffer, binary.BigEndian, ifb.Lat)
		binary.Write(Buffer, binary.BigEndian, ifb.Wins)
		binary.Write(Buffer, binary.BigEndian, ifb.Wind)
		elem = elem.Next()
	}
	binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-answer.Head.NbrOctet))
	buf = Buffer.Bytes()
	fmt.Println(buf)
	return buf
}

/* Manage_type_1 Remplie le buffer avec une reponse au type 1 de la requete avec un maximum de 10 ballons et de 1 packet par requete*/
func Manage_type_1(Req *list.Element, Data *Data) {
	list_tmp := list.New()

	elem := Data.Lst_ball.Lst.Front()
	for elem != nil {
		Coord := elem.Value.(ballon.Ball).Coord.Value.(ballon.Checkpoints).Coord
		if Coord.Longitude < Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Longitude+0.01 && Coord.Longitude > Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Longitude-0.01 && Coord.Latitude < Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Latitude+0.01 && Coord.Latitude > Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Latitude-0.01 {
			ball := elem.Value.(ballon.Ball)
			var ifball InfoBall
			ifball.Id = ball.Id_ball
			ifball.Title = ball.Name
			ifball.Lon = ball.Coord.Value.(ballon.Checkpoints).Coord.Longitude
			ifball.Lat = ball.Coord.Value.(ballon.Checkpoints).Coord.Latitude
			ifball.Wins = ball.Wind.Speed
			ifball.Wind = ball.Wind.Degress
			list_tmp.PushFront(ifball)
		}
		elem = elem.Next()
	}
	Len := list_tmp.Len()
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
	answer := Write_type1(Req, list_tmp)
	Data.Lst_asw.PushBack(answer)
}

/*func Str_cut_n(str string, n int) (s1 string, s2 string) {
	s1 = make(string)
	s2 = ""
	i := 0
	for i < n {
		s1 + 1 = str[i]
		i++
	}
	j := 0
	size := len(str)
	for i < size {
		s2[i] = str[j]
		i++
		j++
	}
	return s1, s2
}*/

func Cut_messagemultipack(pack Packet, msg Message) (tmp_lst *list.List) {
	tmp_lst = list.New()
	size := msg.Size

	for size != 0 {
		if 16+pack.Head.NbrOctet >= 1024 {
			tmp_lst.PushBack(pack)
			pack = Packet{}
			typesp := Pack7{}
			pack.Head.NbrOctet = 24
			typesp.Mess = list.New()
			pack.TPack = typesp
		} else {
			newMess := Message{}
			if msg.Size+16+pack.Head.NbrOctet < 1024 {
				pack.Head.NbrOctet += 16 + msg.Size
				newMess.Size = msg.Size
				newMess.Idmess = msg.Idmess
				newMess.Idcountry = msg.Idcountry
				newMess.Idcity = msg.Idcity
				newMess.Message = msg.Message
				newMess.Type_ = msg.Type_
				pack.TPack.(Pack7).Mess.PushBack(newMess)
				size -= msg.Size
			} else {
				newMess.Size = 1024 - pack.Head.NbrOctet - 16
				size = size - newMess.Size
				pack.Head.NbrOctet += 16 + newMess.Size
				newMess.Idmess = msg.Idmess
				newMess.Idcountry = msg.Idcountry
				newMess.Idcity = msg.Idcity
				//newMess.Message = n octets de la string
				newMess.Type_ = msg.Type_
				pack.TPack.(Pack7).Mess.PushBack(newMess)
			}
		}
	}
	tmp := tmp_lst.Back()
	if tmp.Value.(Packet) != pack {
		tmp_lst.PushBack(pack)
	}
	return tmp_lst
}

func Write_type_2(Ball ballon.Ball) (lst_tmsg *list.List) {
	fmt.Println("Write_type_2")
	var pack Packet
	var typesp Pack7

	lst_pack := list.New()
	pack.Head.NbrOctet = 24
	pack.Head.NumPackReq = 0
	typesp.Mess = list.New()
	pack.TPack = typesp
	elem := Ball.Lst_msg.Front()
	for elem != nil {
		msg := elem.Value.(Message)
		if msg.Size+16+pack.Head.NbrOctet > 1024 {
			tmp_lst := Cut_messagemultipack(pack, msg)
			tmp := tmp_lst.Back()
			pack = (tmp_lst.Remove(tmp)).(Packet)
			lst_pack.PushBackList(tmp_lst)
		} else {
			newMess := Message{}
			pack.Head.NbrOctet += 16 + msg.Size
			newMess.Size = msg.Size
			newMess.Idmess = msg.Idmess
			newMess.Idcountry = msg.Idcountry
			newMess.Idcity = msg.Idcity
			newMess.Message = msg.Message
			newMess.Type_ = msg.Type_
			pack.TPack.(Pack7).Mess.PushBack(newMess)
		}
		elem = elem.Next()
	}
	tmp := lst_pack.Back()
	if tmp.Value.(Packet) != pack {
		lst_pack.PushBack(pack)
	}
	elem = lst_pack.Front()
	numPack := 0
	NbrPack := lst_pack.Len()
	follow := Ball.List_follow.Len()
	for elem != nil {
		tp := elem.Value.(Packet)
		tp.Head.TypeReq = 7
		tp.Head.NbrPackReq = (int32)(NbrPack)
		tp.Head.NumPackReq = (int32)(numPack)
		numPack += 1
		tmp := Pack7{}
		tmp.Mess = tp.TPack.(Pack7).Mess
		tmp.Nbruser = (int32)(follow)
		tmp.NbrMess = (int32)(tmp.Mess.Len())
		elem = elem.Next()
	}
	lst_tmsg = list.New()
	elem = lst_pack.Front()
	for elem != nil {
		tpack := elem.Value.(Packet)
		Buffer := Write_header(tpack)
		binary.Write(Buffer, binary.BigEndian, tpack.TPack.(Pack7).Nbruser)
		binary.Write(Buffer, binary.BigEndian, tpack.TPack.(Pack7).NbrMess)
		tmess := tpack.TPack.(Pack7).Mess.Front()
		for tmess != nil {
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).Idmess)
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).Size)
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).Idcountry)
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).Idcity)
			Buffer.WriteString(tmess.Value.(Message).Message)
			tmess = tmess.Next()
		}
		elem = elem.Next()
		lst_tmsg.PushBack(Buffer.Bytes())
	}
	return lst_tmsg
}

/* Manage_type_2 Remplie le buffer avec une reponse au type 2 de la requete. */
func Manage_type_2(Req *list.Element, Data *Data) {
	elem := Data.Lst_ball.Lst.Front()
	for elem != nil && elem.Value.(ballon.Ball).Id_ball != Req.Value.(protocol.Lst_req_sock).Union.(protocol.Id_ballon).IdBallon {
		elem = elem.Next()
	}
	if elem != nil {
		ball := elem.Value.(ballon.Ball)
		ball.Possessed = &Data.User
		ball.List_follow.PushFront(Data.User)
		Lst_answer := Write_type_2(ball)
		Data.Lst_asw.PushBackList(Lst_answer)
	} else {
		// Acknowledgement negatif
		fmt.Println("Ballon not found") // Ball not found
	}
}

/* Manage_type_3 Remplie le buffer avec une reponse au type 3 de la requete. */
func Manage_type_3(Req *list.Element, Data *Data) {

}

/* Manage_type_4 Remplie le buffer avec une reponse au type 4 de la requete. */
func Manage_type_4(Req *list.Element, Data *Data) {

}

/* Manage_type_5 Remplie le buffer avec une reponse au type 5 de la requete. */
func Manage_type_5(Req *list.Element, Data *Data) {

}

/*
** Get_answer fournis une reponse approprie a la requete du client,
** avec un buffer de 1024 Octets. Elle initialisera l'authentification de
** de l'utilisateur et nettoiera le flux de requetes traites.
 */
func (Data *Data) Get_answer() (err error) {
	fmt.Println("Get_answer")
	Req := Data.Lst_req.Front()
	fmt.Println(Req.Value.(protocol.Lst_req_sock))
	if Req == nil {
		fmt.Println("Error get answer")
	}
	fmt.Println(Data)
	Data.User, err = Check_user(Req, Data.Lst_users)
	if err != nil {
		fmt.Println(err)
	} else {
		if Check_packets_list(Req) == true {
			switch Req.Value.(protocol.Lst_req_sock).Type {
			case 1:
				Manage_type_1(Req, Data)
			case 2:
				Manage_type_2(Req, Data)
			case 3:
				Manage_type_3(Req, Data)
			case 4:
				Manage_type_4(Req, Data)
			case 5:
				Manage_type_5(Req, Data)
			}
		}
		Del_request_done(Data.Lst_req)
		return nil
	}
	Del_request_done(Data.Lst_req)
	err = errors.New("Failed check user")
	return err
}

/*
** Get_acknowledgement creer un buffer pour confirmation au client et le rempli.
 */
func Get_aknowledgement(Lst_req *list.List, Lst_usr *users.All_users) (answer []byte, err error) {
	return answer, err
}
