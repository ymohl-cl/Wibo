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
	"owm"
	"protocol"
	"time"
	"users"
)

const (
	_ = iota
	// SAME FOR CLIENT AND SERVER
	ACK   = 32767
	TAKEN = 4
	// DEFINE CLIENT
	SYNC       = 1
	MAJ        = 2
	POS        = 3
	FOLLOW_ON  = 5
	FOLLOW_OFF = 6
	NEW_BALL   = 7
	SEND_BALL  = 8
	// DEFINE SERVER
	CONN     = 1
	INF_BALL = 2
	NEARBY   = 3
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
	Message   string // []byte
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

type Mini_ball struct {
	Id_ball  int64
	Nbr_mess int32
	Flag     int32
}

type Pack1 struct {
	Nbr_ball   int16
	mini_balls *list.List
}

type Header struct {
	NbrOctet   int16
	TypeReq    int16
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
func (Data *Data) Check_lstrequest() bool {
	Req := Data.Lst_req.Front()
	next := Req.Next()
	tmp := Req
	var nr protocol.Lst_req_sock
	var tr protocol.Lst_req_sock

	tr = tmp.Value.(protocol.Lst_req_sock)
	fmt.Println("Value dans check de tr:")
	fmt.Println(tr)
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

func Write_type3(Req *list.Element, list_tmp *list.List) (buf []byte) {
	var answer Packet
	var typesp Pack6

	answer.Head.NbrOctet = 24
	answer.Head.TypeReq = POS
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

func Str_cut_n(str string, n int) (s1 string, s2 string) {
	Bstr := bytes.NewBufferString(str)
	Bs1 := Bstr.Next(n)
	s2 = Bstr.String()
	tmp := bytes.NewBuffer(Bs1)
	s1 = tmp.String()
	return s1, s2
}

func Cut_messagemultipack(pack Packet, msg1 ballon.Lst_msg) (tmp_lst *list.List) {
	tmp_lst = list.New()
	size := msg1.Size
	var msg ballon.Lst_msg

	msg = msg1
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
			if int16(msg.Size)+16+pack.Head.NbrOctet < 1024 {
				pack.Head.NbrOctet += 16 + (int16)(msg.Size)
				newMess.Size = msg.Size
				newMess.Idmess = msg.Id_Message
				newMess.Idcountry = msg.Idcountry
				newMess.Idcity = msg.Idcity
				newMess.Message = msg.Content
				newMess.Type_ = msg.Type_
				pack.TPack.(Pack7).Mess.PushBack(newMess)
				size -= msg.Size
			} else {
				newMess.Size = 1024 - int32(pack.Head.NbrOctet) - 16
				size = size - newMess.Size
				pack.Head.NbrOctet += 16 + int16(newMess.Size)
				newMess.Idmess = msg.Id_Message
				newMess.Idcountry = msg.Idcountry
				newMess.Idcity = msg.Idcity
				newMess.Message, msg.Content = Str_cut_n(msg.Content, int(newMess.Size))
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

func Write_type_4(Ball ballon.Ball) (lst_tmsg *list.List) {
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
		fmt.Println("Mess of ball found")
		msg := elem.Value.(ballon.Lst_msg)
		if int16(msg.Size)+16+pack.Head.NbrOctet > 1024 {
			fmt.Println("Depassement de buffer !cut! ")
			tmp_lst := Cut_messagemultipack(pack, msg)
			tmp := tmp_lst.Back()
			pack = (tmp_lst.Remove(tmp)).(Packet)
			lst_pack.PushBackList(tmp_lst)
		} else {
			fmt.Println("In the buffer OK")
			newMess := Message{}
			pack.Head.NbrOctet += 16 + int16(msg.Size)
			newMess.Size = msg.Size
			newMess.Idmess = msg.Id_Message
			newMess.Idcountry = msg.Idcountry
			newMess.Idcity = msg.Idcity
			newMess.Message = msg.Content
			newMess.Type_ = msg.Type_
			pack.TPack.(Pack7).Mess.PushBack(newMess)
		}
		elem = elem.Next()
	}
	tmp := lst_pack.Back()
	if tmp == nil || tmp.Value.(Packet) != pack {
		lst_pack.PushBack(pack)
	}
	elem = lst_pack.Front()
	numPack := 0
	NbrPack := lst_pack.Len()
	if NbrPack == 0 {
		NbrPack = 1
	}
	follow := Ball.List_follow.Len()
	for elem != nil {
		tp := elem.Value.(Packet)
		tp.Head.TypeReq = TAKEN
		tp.Head.NbrPackReq = (int32)(NbrPack)
		tp.Head.NumPackReq = (int32)(numPack)
		numPack += 1
		tmp := Pack7{}
		tmp.Mess = tp.TPack.(Pack7).Mess
		tmp.Nbruser = (int32)(follow - 1) // - l'ajout qui viens d'etre fait.
		tmp.NbrMess = (int32)(tmp.Mess.Len())
		tp.TPack = tmp
		elem.Value = tp
		elem = elem.Next()
	}
	lst_tmsg = list.New()
	elem = lst_pack.Front()
	for elem != nil {
		tpack := elem.Value.(Packet)
		fmt.Println("Debeug")
		fmt.Println(tpack)
		fmt.Println(elem.Value.(Packet))
		Buffer := Write_header(tpack)
		binary.Write(Buffer, binary.BigEndian, Ball.Id_ball)
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
		binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-tpack.Head.NbrOctet))
		elem = elem.Next()
		lst_tmsg.PushBack(Buffer.Bytes())
	}
	return lst_tmsg
}

/* Write_type_Ack */
func Manage_typeack(Type int16, IdMobile int64, IdBallon int64, value int32) (answer []byte) {
	tpack := Packet{}

	tpack.Head.NbrOctet = int16(72)
	tpack.Head.TypeReq = ACK
	tpack.Head.NbrPackReq = int32(1)
	tpack.Head.NumPackReq = int32(0)
	Buffer := Write_header(tpack)
	binary.Write(Buffer, binary.BigEndian, Type)
	binary.Write(Buffer, binary.BigEndian, value)
	if Type == NEW_BALL {
		binary.Write(Buffer, binary.BigEndian, int64(0))
	} else {
		binary.Write(Buffer, binary.BigEndian, IdBallon)
	}
	//	binary.Write(Buffer, binary.BigEndian, IdMobile)
	//	binary.Write(Buffer, binary.BigEndian, make([]byte, 32))
	binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-tpack.Head.NbrOctet))
	answer = Buffer.Bytes()
	return answer
}

func Write_type_1(lst_tmp *list.List) (asw *list.List) {
	asw = list.New()

	elem := lst_tmp.Front()
	for elem != nil {
		tpack := elem.Value.(Packet)
		Buffer := Write_header(tpack)
		binary.Write(Buffer, binary.BigEndian, tpack.TPack.(Pack1).Nbr_ball)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
		telem := tpack.TPack.(Pack1).mini_balls.Front()
		for telem != nil {
			mball := telem.Value.(Mini_ball)
			binary.Write(Buffer, binary.BigEndian, mball.Id_ball)
			binary.Write(Buffer, binary.BigEndian, mball.Nbr_mess)
			binary.Write(Buffer, binary.BigEndian, mball.Flag)
			telem = telem.Next()
		}
		asw.PushBack(Buffer.Bytes())
		elem = elem.Next()
	}
	return (asw)
}

/* Manage_type_1 Remplie le buffer avec une reponse au type 1 de la requete avec tous les ballons possedes et suivis par l'utilisateur.*/
func (Data *Data) Manage_sync(Req *list.Element) {
	elem := Data.User.List_follow.Front()
	lst_tmp := list.New()
	var pack Packet
	var pack1 Pack1

	pack.Head.NbrOctet = 24
	pack.Head.TypeReq = CONN
	pack1.Nbr_ball = 0
	pack1.mini_balls = list.New()
	for elem != nil {
		var mball Mini_ball
		ball := elem.Value.(ballon.Ball)
		mball.Id_ball = ball.Id_ball
		mball.Nbr_mess = int32(ball.Lst_msg.Len())
		if ball.Possessed == &Data.User {
			mball.Flag = 1
		} else {
			mball.Flag = 0
		}
		if pack.Head.NbrOctet+16 > 1024 {
			pack.TPack = pack1
			lst_tmp.PushBack(pack)
			pack = Packet{}
			pack1 = Pack1{}
			pack.Head.NbrOctet = 24
			pack.Head.TypeReq = CONN
			pack1.Nbr_ball = 0
			pack1.mini_balls = list.New()
		}
		pack.Head.NbrOctet += 16
		pack1.mini_balls.PushBack(mball)
		pack1.Nbr_ball += 1
		elem = elem.Next()
	}
	if lst_tmp.Back().Value.(Packet) != pack {
		lst_tmp.PushBack(pack)
	}
	elem = lst_tmp.Front()
	nbrPack := int32(lst_tmp.Len())
	numPack := int32(0)
	for elem != nil {
		packet := elem.Value.(Packet)
		packet.Head.NbrPackReq = nbrPack
		packet.Head.NumPackReq = numPack
		numPack += 1
		elem = elem.Next()
	}
	Data.Lst_asw.PushBackList(Write_type_1(lst_tmp))
}

/* Fournis le contenu d'un ballon au client, reponse au type 2 par un type 2 */
func (Data *Data) Manage_maj(Req *list.Element) {
}

func (Data *Data) Manage_sendball(Req *list.Element) {
}

/* Manage_type_3 Remplie le buffer avec une reponse au type 3 de la requete avec un maximum de 10 ballons et de 1 packet par requete*/
func (Data *Data) Manage_pos(Req *list.Element) {
	list_tmp := list.New()
	var ifball InfoBall

	elem := Data.Lst_ball.Lst.Front()
	for elem != nil {
		Coord := elem.Value.(ballon.Ball).Coord.Value.(ballon.Checkpoints).Coord
		fmt.Println("Balon testing and coordinates:")
		fmt.Println(Coord.Longitude)
		fmt.Println(Coord.Latitude)
		fmt.Println("Position testing and coordinates:")
		fmt.Println(Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Longitude + 0.01)
		fmt.Println(Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Latitude + 0.01)

		if Coord.Longitude < Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Longitude+0.01 && Coord.Longitude > Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Longitude-0.01 && Coord.Latitude < Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Latitude+0.01 && Coord.Latitude > Req.Value.(protocol.Lst_req_sock).Union.(protocol.Position).Latitude-0.01 {
			ball := elem.Value.(ballon.Ball)
			ifball.Id = ball.Id_ball
			ifball.Title = ball.Name
			ifball.Lon = ball.Coord.Value.(ballon.Checkpoints).Coord.Longitude
			ifball.Lat = ball.Coord.Value.(ballon.Checkpoints).Coord.Latitude
			fmt.Println("Balon testing:")
			fmt.Println(Coord)

			ifball.Wins = ball.Wind.Speed
			ifball.Wind = ball.Wind.Degress
			list_tmp.PushBack(ifball)
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
	answer := Write_type3(Req, list_tmp)
	Data.Lst_asw.PushBack(answer)
}

/* Manage_type_4 Remplie le buffer avec une reponse au type 2 ou 4 de la requete. */
func (Data *Data) Manage_taken(Req *list.Element) {
	elem := Data.Lst_ball.Lst.Front()
	treq := Req.Value.(protocol.Lst_req_sock)

	fmt.Println("Debug Manage_taken")
	fmt.Println("Nombre de ballon: ")
	fmt.Println(Data.Lst_ball.Lst.Front)

	for elem != nil && elem.Value.(ballon.Ball).Id_ball != Req.Value.(protocol.Lst_req_sock).Union.(protocol.Id_ballon).IdBallon {
		elem = elem.Next()
	}
	if elem != nil {
		fmt.Println("Ball found")
		fmt.Println(elem.Value.(ballon.Ball).Id_ball)
		fmt.Println(Req.Value.(protocol.Lst_req_sock).Union.(protocol.Id_ballon).IdBallon)
		ball := elem.Value.(ballon.Ball)
		ball.Possessed = &Data.User
		// Rajouter checkpoint actuel au ballon qui est celui de l'utilisateur.
		// Dans l'application move checkpoint, voir pour que si le ballon est possessed il ne faut pas qu'il se deplace.
		// Verifier si le ballon n'est pas deja possessed
		ball.List_follow.PushFront(Data.User)
		Lst_answer := Write_type_4(ball)
		Data.Lst_asw.PushBackList(Lst_answer)
	} else {
		answer := Manage_typeack(treq.Type, treq.IdMobile, treq.Union.(protocol.Id_ballon).IdBallon, int32(0))
		Data.Lst_asw.PushBack(answer)
		fmt.Println("Ballon not found")
	}
}

/* Manage_type_3 Remplie le buffer avec une reponse au type 3 de la requete. */
func (Data *Data) Manage_followon(Req *list.Element) {
	treq := Req.Value.(protocol.Lst_req_sock)

	elem := Data.Lst_ball.Lst.Front()
	for elem != nil && elem.Value.(ballon.Ball).Id_ball != treq.Union.(protocol.Id_ballon).IdBallon {
		elem = elem.Next()
	}
	var answer []byte
	if elem != nil {
		answer = Manage_typeack(treq.Type, treq.IdMobile, treq.Union.(protocol.Id_ballon).IdBallon, int32(0))
	} else {
		elem.Value.(ballon.Ball).List_follow.PushBack(Data.User)
		answer = Manage_typeack(treq.Type, treq.IdMobile, treq.Union.(protocol.Id_ballon).IdBallon, int32(1))
	}
	Data.Lst_asw.PushBack(answer)
}

/* Manage_type_4 Remplie le buffer avec une reponse au type 4 de la requete. */
func (Data *Data) Manage_followoff(Req *list.Element) {
	treq := Req.Value.(protocol.Lst_req_sock)
	elem := Data.Lst_ball.Lst.Front()
	for elem != nil && elem.Value.(ballon.Ball).Id_ball != treq.Union.(protocol.Id_ballon).IdBallon {
		elem = elem.Next()
	}
	var answer []byte
	var dvc *list.Element = nil
	var users_l *list.Element = nil
	if elem != nil {
		users_l = elem.Value.(ballon.Ball).List_follow.Front()
		for users_l != nil {
			user := users_l.Value.(users.User)
			dvc = user.Device.Front()
			for dvc != nil && dvc.Value.(users.Device).IdMobile != treq.IdMobile {
				dvc = dvc.Next()
			}
			users_l = users_l.Next()
		}
	}
	if dvc != nil {
		elem.Value.(ballon.Ball).List_follow.Remove(users_l)
		answer = Manage_typeack(treq.Type, treq.IdMobile, treq.Union.(protocol.Id_ballon).IdBallon, int32(1))
	} else {
		answer = Manage_typeack(treq.Type, treq.IdMobile, treq.Union.(protocol.Id_ballon).IdBallon, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
}

/* Manage_type_5 Remplie le buffer avec une reponse au type 5 de la requete. */
func (Data *Data) Manage_newball(Req *list.Element, Tab_wd *owm.All_data) {
	var ball ballon.Ball
	var newBall protocol.Ballon
	var mess ballon.Lst_msg

	newBall = Req.Value.(protocol.Lst_req_sock).Union.(protocol.Ballon)
	Data.Lst_ball.Id_max++
	ball.Id_ball = Data.Lst_ball.Id_max
	ball.Name = newBall.Title
	ball.Lst_msg = list.New()
	mess.Id_Message = 0
	mess.Size = (int32)(len(newBall.Message))
	mess.Content = newBall.Message
	mess.Type_ = 1
	ball.Lst_msg.PushFront(mess)
	ball.Date = time.Now()
	ball.Possessed = &(Data.User)
	ball.List_follow.PushFront(Data.User)
	ball.Creator = &(Data.User)
	ball.Get_checkpointList(Tab_wd.Get_Paris())
	Data.Lst_ball.Lst.PushFront(ball)
	treq := Req.Value.(protocol.Lst_req_sock)
	answer := Manage_typeack(treq.Type, treq.IdMobile, ball.Id_ball, int32(1))
	Data.Lst_asw.PushBack(answer)
}

/*
** Get_answer fournis une reponse approprie a la requete du client,
** avec un buffer de 1024 Octets. Elle initialisera l'authentification de
** de l'utilisateur et nettoiera le flux de requetes traites.
 */
func (Data *Data) Get_answer(Tab_wd *owm.All_data) (err error) {
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
		switch Req.Value.(protocol.Lst_req_sock).Type {
		case 1:
			Data.Manage_sync(Req)
		case 2:
			Data.Manage_maj(Req)
		case 3:
			Data.Manage_pos(Req)
		case 4:
			Data.Manage_taken(Req)
		case 5:
			Data.Manage_followon(Req)
		case 6:
			Data.Manage_followoff(Req)
		case 7:
			Data.Manage_newball(Req, Tab_wd)
		case 8:
			Data.Manage_sendball(Req)
		case ACK:
			fmt.Println("Acknowldgement received")
		}
		Del_request_done(Data.Lst_req)
		return nil
	}
	Del_request_done(Data.Lst_req)
	err = errors.New("Failed check user")
	return err
}

/*
** Create a new
 */
func (Data *Data) Get_aknowledgement(Lst_usr *users.All_users) (answer []byte) {
	elem := Data.Lst_req.Back()
	treq := elem.Value.(protocol.Lst_req_sock)

	if treq.Type == 7 {
		answer = Manage_typeack(treq.Type, treq.IdMobile, 0, int32(1))
	} else {
		answer = Manage_typeack(treq.Type, treq.IdMobile, treq.Union.(protocol.Id_ballon).IdBallon, int32(1))
	}
	return answer
}
