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
	// SAME DEFINE FOR CLIENT AND SERVER
	ACK   = 32767
	TAKEN = 4
	// DEFINE CLIENT
	SYNC       = 1
	UPDATE     = 2
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

type Posball struct {
	id    int64
	title string
	lon   float64
	lat   float64
	wins  float64
	wind  float64
}

type Nearby struct {
	nbrball int32
	balls   *list.List
}

type Message struct {
	id        int32
	size      int32
	idcountry int32
	idcity    int32
	mess      string // []byte
	mtype     int32
}

type Contentball struct {
	nbruser  int32
	nbrmess  int32
	messages *list.List
}

type Infoball struct {
	id      int64
	nbrmess int32
	taken   int32
}

type Conn struct {
	nbrball   int32
	infoballs *list.List
}

type Header struct {
	octets int16
	rtype  int16
	pnbr   int32
	pnum   int32
}

type Packet struct {
	head  Header
	ptype interface{}
}

type Data struct {
	Lst_req   *list.List
	Lst_asw   *list.List
	Lst_ball  *ballon.All_ball
	Lst_users *users.All_users
	User      *list.Element
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
		if 16+pack.head.octets >= 1024 {
			tmp_lst.PushBack(pack)
			pack = Packet{}
			typesp := Contentball{}
			pack.head.octets = 24
			typesp.messages = list.New()
			pack.ptype = typesp
		} else {
			newMess := Message{}
			if int16(msg.Size)+16+pack.head.octets < 1024 {
				pack.head.octets += 16 + (int16)(msg.Size)
				newMess.size = msg.Size
				newMess.id = msg.Id_Message
				newMess.idcountry = msg.Idcountry
				newMess.idcity = msg.Idcity
				newMess.mess = msg.Content
				newMess.mtype = msg.Type_
				pack.ptype.(Contentball).messages.PushBack(newMess)
				size -= msg.Size
			} else {
				newMess.size = 1024 - int32(pack.head.octets) - 16
				size = size - newMess.size
				pack.head.octets += 16 + int16(newMess.size)
				newMess.id = msg.Id_Message
				newMess.idcountry = msg.Idcountry
				newMess.idcity = msg.Idcity
				newMess.mess, msg.Content = Str_cut_n(msg.Content, int(newMess.size))
				newMess.mtype = msg.Type_
				pack.ptype.(Contentball).messages.PushBack(newMess)
			}
		}
	}
	tmp := tmp_lst.Back()
	if tmp.Value.(Packet) != pack {
		tmp_lst.PushBack(pack)
	}
	return tmp_lst
}

/*
** Check_packets_list Verifie la validite des packets suivants en partant du
** principe que le header du premier packet est valide (traitement multi-packet).
 */
func (Data *Data) Check_lstrequest() bool {
	Req := Data.Lst_req.Front()
	next := Req.Next()
	tmp := Req
	var nr protocol.Request
	var tr protocol.Request

	tr = tmp.Value.(protocol.Request)
	fmt.Println("Value dans check de tr:")
	fmt.Println(tr)
	for next != nil {
		nr = next.Value.(protocol.Request)
		if tr.Nbrpck == tr.Numpck+1 {
			return true
		} else if tr.Rtype != nr.Rtype {
			return false
		} else if tr.Numpck != nr.Numpck+1 {
			return false
		}
		next = next.Next()
		tmp = tmp.Next()
		tr = tmp.Value.(protocol.Request)
	}
	if tr.Nbrpck == tr.Numpck+1 {
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
		if elem.Value.(protocol.Request).Numpck == elem.Value.(protocol.Request).Nbrpck-1 {
			Lst_req.Remove(elem)
			return
		}
		tmp := elem
		Lst_req.Remove(tmp)
		elem = elem.Next()
	}
}

func Write_header(answer Packet) (Buffer *bytes.Buffer) {
	Buffer = new(bytes.Buffer)

	fmt.Println("Header")
	fmt.Println(answer)
	binary.Write(Buffer, binary.BigEndian, answer.head.octets)
	binary.Write(Buffer, binary.BigEndian, answer.head.rtype)
	binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
	binary.Write(Buffer, binary.BigEndian, answer.head.pnbr)
	binary.Write(Buffer, binary.BigEndian, answer.head.pnum)
	return Buffer
}

/* Write a type connexion, return answer list */
func Write_conn(plist *list.List) (alist *list.List) {
	alist = list.New()

	packet := plist.Front()
	for packet != nil {
		pck := packet.Value.(*Packet)
		Buffer := Write_header(*pck)
		binary.Write(Buffer, binary.BigEndian, pck.ptype.(Conn).nbrball)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
		iball := pck.ptype.(Conn).infoballs.Front()
		for iball != nil {
			ball := iball.Value.(Infoball)
			binary.Write(Buffer, binary.BigEndian, ball.id)
			binary.Write(Buffer, binary.BigEndian, ball.nbrmess)
			binary.Write(Buffer, binary.BigEndian, ball.taken)
			iball = iball.Next()
		}
		binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-pck.head.octets))
		alist.PushBack(Buffer.Bytes())
		packet = packet.Next()
	}
	return (alist)
}

/* Write_type_Ack */
func Manage_ack(Type int16, IdMobile int64, IdBallon int64, value int32) (answer []byte) {
	tpack := Packet{}

	tpack.head.octets = int16(32)
	tpack.head.rtype = ACK
	tpack.head.pnbr = int32(1)
	tpack.head.pnum = int32(0)
	Buffer := Write_header(tpack)
	binary.Write(Buffer, binary.BigEndian, int32(Type))
	binary.Write(Buffer, binary.BigEndian, value)
	binary.Write(Buffer, binary.BigEndian, IdBallon)
	binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-tpack.head.octets))
	answer = Buffer.Bytes()
	return answer
}

func Write_nearby(Req *list.Element, list_tmp *list.List) (buf []byte) {
	var answer Packet
	var typesp Nearby

	answer.head.octets = 24
	answer.head.rtype = POS
	answer.head.pnbr = 1
	answer.head.pnum = 0
	typesp.nbrball = (int32)(list_tmp.Len())
	typesp.balls = list_tmp
	answer.ptype = typesp
	elem := list_tmp.Front()
	for elem != nil {
		answer.head.octets += 56
		elem = elem.Next()
	}
	Buffer := Write_header(answer)
	binary.Write(Buffer, binary.BigEndian, answer.ptype.(Nearby).nbrball)
	binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
	elem = answer.ptype.(Nearby).balls.Front()
	for elem != nil {
		ifb := elem.Value.(Posball)
		binary.Write(Buffer, binary.BigEndian, ifb.id)
		Buffer.WriteString(ifb.title)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 16-len(ifb.title)))
		binary.Write(Buffer, binary.BigEndian, ifb.lon)
		binary.Write(Buffer, binary.BigEndian, ifb.lat)
		binary.Write(Buffer, binary.BigEndian, ifb.wins)
		binary.Write(Buffer, binary.BigEndian, ifb.wind)
		elem = elem.Next()
	}
	binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-answer.head.octets))
	buf = Buffer.Bytes()
	fmt.Println(buf)
	return buf
}

func Write_contentball(Ball *ballon.Ball, packettype int16) (alist *list.List) {
	var pack Packet
	var contball Contentball

	alist = list.New()
	plist := list.New()
	pack.head.octets = 32
	pack.head.pnum = 0
	contball.messages = list.New()
	pack.ptype = contball
	emess := Ball.Lst_msg.Front()
	for emess != nil {
		msg := emess.Value.(ballon.Lst_msg)
		if int16(msg.Size)+16+pack.head.octets > 1024 {
			tmp_lst := Cut_messagemultipack(pack, msg)
			tmp := tmp_lst.Back()
			pack = (tmp_lst.Remove(tmp)).(Packet)
			plist.PushBackList(tmp_lst)
		} else {
			newMess := Message{}
			pack.head.octets += 16 + int16(msg.Size)
			newMess.size = msg.Size
			newMess.id = msg.Id_Message
			newMess.idcountry = msg.Idcountry
			newMess.idcity = msg.Idcity
			newMess.mess = msg.Content
			newMess.mtype = msg.Type_
			pack.ptype.(Contentball).messages.PushBack(newMess)
		}
		emess = emess.Next()
	}
	tmp := plist.Back()
	if tmp == nil || tmp.Value.(Packet) != pack {
		plist.PushBack(pack)
	}
	epck := plist.Front()
	numPack := 0
	NbrPack := plist.Len()
	if NbrPack == 0 {
		NbrPack = 1
	}
	follow := Ball.List_follow.Len()
	for epck != nil {
		tp := epck.Value.(Packet)
		tp.head.rtype = packettype
		tp.head.pnbr = (int32)(NbrPack)
		tp.head.pnum = (int32)(numPack)
		numPack += 1
		tmp := Contentball{}
		tmp.messages = tp.ptype.(Contentball).messages
		tmp.nbruser = (int32)(follow - 1) // - l'ajout qui viens d'etre fait.
		tmp.nbrmess = (int32)(tmp.messages.Len())
		tp.ptype = tmp
		epck.Value = tp
		epck = epck.Next()
	}
	epck = plist.Front()
	for epck != nil {
		tpack := epck.Value.(Packet)
		Buffer := Write_header(tpack)
		binary.Write(Buffer, binary.BigEndian, Ball.Id_ball)
		binary.Write(Buffer, binary.BigEndian, tpack.ptype.(Contentball).nbruser)
		binary.Write(Buffer, binary.BigEndian, tpack.ptype.(Contentball).nbrmess)
		tmess := tpack.ptype.(Contentball).messages.Front()
		for tmess != nil {
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).id)
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).size)
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).idcountry)
			binary.Write(Buffer, binary.BigEndian, tmess.Value.(Message).idcity)
			Buffer.WriteString(tmess.Value.(Message).mess)
			tmess = tmess.Next()
		}
		binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-tpack.head.octets))
		epck = epck.Next()
		alist.PushBack(Buffer.Bytes())
	}
	return alist
}

/* Get all follower ball and create packets list for send the infoballs */
func (Data *Data) Manage_sync(Req *list.Element) {
	fmt.Println("Manage_sync")
	flwupball := Data.User.Value.(*users.User).List_follow.Front()
	plist := list.New()
	pck := new(Packet)
	conn := new(Conn)

	pck.head.octets = 24
	pck.head.rtype = CONN
	conn.nbrball = 0
	conn.infoballs = list.New()

	for flwupball != nil {
		myball := new(Infoball)
		ball := flwupball.Value.(*list.Element).Value.(*ballon.Ball)
		myball.id = ball.Id_ball
		myball.nbrmess = int32(ball.Lst_msg.Len())
		if ball.Possessed == Data.User {
			myball.taken = 1
		} else {
			myball.taken = 0
		}
		if pck.head.octets+16 > 1024 {
			pck.ptype = conn
			plist.PushBack(pck)
			pck = new(Packet)
			conn = new(Conn)
			pck.head.octets = 24
			pck.head.rtype = CONN
			conn.nbrball = 0
			conn.infoballs = list.New()
		}
		pck.head.octets += 16
		conn.infoballs.PushBack(*myball)
		conn.nbrball += 1
		flwupball = flwupball.Next()
	}
	pck.ptype = *conn
	plist.PushBack(pck)
	pelem := plist.Front()
	nbrpacket := int32(plist.Len())
	numpacket := int32(0)
	fmt.Println("info des ballons followed nbrpacker numpacket")
	fmt.Println(nbrpacket)
	fmt.Println(numpacket)
	for pelem != nil {
		fmt.Println("Coucou")
		packet := pelem.Value.(*Packet)
		packet.head.pnbr = nbrpacket
		packet.head.pnum = numpacket
		numpacket += 1
		pelem = pelem.Next()
	}
	Data.Lst_asw.PushBackList(Write_conn(plist))
}

/* If user has idbaloon in follower list, give him data message */
func (Data *Data) Manage_update(request *list.Element) {
	var ball *ballon.Ball

	rqt := request.Value.(protocol.Request)
	idsearch := rqt.Spec.(protocol.Ballid).Id
	eball := Data.User.Value.(*users.User).List_follow.Front()
	for eball != nil {
		ball = eball.Value.(*list.Element).Value.(*ballon.Ball)
		if ball.Id_ball == idsearch {
			break
		}
		eball = eball.Next()
	}
	if eball != nil {
		Lst_answer := Write_contentball(ball, UPDATE)
		Data.Lst_asw.PushBackList(Lst_answer)
	} else {
		answer := Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Ballid).Id, int32(0))
		Data.Lst_asw.PushBack(answer)
	}
}

/* Manage_type_3 Remplie le buffer avec une reponse au type 3 de la requete avec un maximum de 10 ballons et de 1 packet par requete*/
func (Data *Data) Manage_pos(Req *list.Element) {
	list_tmp := list.New()
	var ifball Posball

	elem := Data.Lst_ball.Lst.Front()
	for elem != nil {
		Coord := elem.Value.(*ballon.Ball).Coord.Value.(ballon.Checkpoints).Coord
		if Coord.Longitude < Req.Value.(protocol.Request).Spec.(protocol.Position).Lon+0.01 && Coord.Longitude > Req.Value.(protocol.Request).Spec.(protocol.Position).Lon-0.01 && Coord.Latitude < Req.Value.(protocol.Request).Spec.(protocol.Position).Lat+0.01 && Coord.Latitude > Req.Value.(protocol.Request).Spec.(protocol.Position).Lat-0.01 {
			ball := elem.Value.(*ballon.Ball)
			ifball.id = ball.Id_ball
			ifball.title = ball.Name
			ifball.lon = ball.Coord.Value.(ballon.Checkpoints).Coord.Longitude
			ifball.lat = ball.Coord.Value.(ballon.Checkpoints).Coord.Latitude
			fmt.Println("Balon testing:")
			fmt.Println(Coord)

			ifball.wins = ball.Wind.Speed
			ifball.wind = ball.Wind.Degress
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
	answer := Write_nearby(Req, list_tmp)
	Data.Lst_asw.PushBack(answer)
}

/* Manage_type_4 Remplie le buffer avec une reponse au type 2 ou 4 de la requete. */
func (Data *Data) Manage_taken(request *list.Element) {
	eball := Data.Lst_ball.Lst.Front()
	rqt := request.Value.(protocol.Request)

	for eball != nil && eball.Value.(*ballon.Ball).Id_ball != request.Value.(protocol.Request).Spec.(protocol.Ballid).Id {
		eball = eball.Next()
	}
	if eball != nil {
		ball := eball.Value.(*ballon.Ball)
		fmt.Println("Print status du ballon")
		fmt.Println(ball.Id_ball)
		fmt.Println(ball.Possessed)
		fmt.Println(ball.Check_userfollower(Data.User))
		if ball.Possessed == nil && ball.Check_userfollower(Data.User) == false {
			ball.Possessed = Data.User
			ball.List_follow.PushFront(Data.User)
			Data.User.Value.(*users.User).List_follow.PushBack(eball)
			Lst_answer := Write_contentball(ball, TAKEN)
			Data.Lst_asw.PushBackList(Lst_answer)
			ball.Clearcheckpoint()
		} else {
			answer := Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Ballid).Id, int32(0))
			Data.Lst_asw.PushBack(answer)
		}
	} else {
		answer := Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Ballid).Id, int32(0))
		Data.Lst_asw.PushBack(answer)
	}
}

/* Check user follower and if is missing, add follower status  */
func (Data *Data) Manage_followon(request *list.Element) {
	rqt := request.Value.(protocol.Request)

	eball := Data.Lst_ball.Lst.Front()
	for eball != nil && eball.Value.(*ballon.Ball).Id_ball != rqt.Spec.(protocol.Ballid).Id {
		eball = eball.Next()
	}
	var answer []byte
	if eball != nil && eball.Value.(*ballon.Ball).Check_userfollower(Data.User) == false {
		answer = Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Ballid).Id, int32(1))
		eball.Value.(*ballon.Ball).List_follow.PushBack(Data.User)
		Data.User.Value.(*users.User).List_follow.PushBack(eball)
	} else {
		answer = Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Ballid).Id, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
}

/* Check user follower and if is found, remove follower status  */
func (Data *Data) Manage_followoff(request *list.Element) {
	var answer []byte
	rqt := request.Value.(protocol.Request)
	eball := Data.Lst_ball.Lst.Front()

	for eball != nil && eball.Value.(*ballon.Ball).Id_ball != rqt.Spec.(protocol.Ballid).Id {
		eball = eball.Next()
	}
	if eball != nil && eball.Value.(*ballon.Ball).Check_userfollower(Data.User) == true {
		eball.Value.(*ballon.Ball).List_follow.Remove(Data.User)
		Data.User.Value.(*users.User).List_follow.Remove(eball)
		answer = Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Ballid).Id, int32(1))

	} else {
		answer = Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Ballid).Id, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
}

/* Manage_type_5 Remplie le buffer avec une reponse au type 5 de la requete. */
func (Data *Data) Manage_newball(requete *list.Element, Tab_wd *owm.All_data) {
	ball := new(ballon.Ball)
	rqt := requete.Value.(protocol.Request)
	var checkpoint ballon.Checkpoints
	var newball protocol.New_ball
	var mess ballon.Lst_msg

	newball = requete.Value.(protocol.Request).Spec.(protocol.New_ball)
	Data.Lst_ball.Id_max++
	ball.Id_ball = Data.Lst_ball.Id_max
	ball.Name = newball.Title
	ball.Lst_msg = list.New()
	ball.List_follow = list.New()
	ball.Checkpoints = list.New()
	mess.Id_Message = 0
	mess.Size = (int32)(len(newball.Message))
	mess.Content = newball.Message
	mess.Type_ = 1
	ball.Lst_msg.PushFront(mess)
	ball.Date = time.Now()
	ball.Possessed = nil
	ball.List_follow.PushFront(Data.User)
	ball.Creator = Data.User
	eball := Data.Lst_ball.Lst.PushBack(ball)
	checkpoint.Coord.Longitude = rqt.Spec.(protocol.New_ball).Lonuser
	checkpoint.Coord.Latitude = rqt.Spec.(protocol.New_ball).Latuser
	eball.Value.(*ballon.Ball).Coord = eball.Value.(*ballon.Ball).Checkpoints.PushBack(checkpoint)
	eball.Value.(*ballon.Ball).Get_checkpointList(Tab_wd.Get_Paris())
	Data.User.Value.(*users.User).List_follow.PushBack(eball)

	answer := Manage_ack(rqt.Rtype, rqt.Deviceid, ball.Id_ball, int32(1))
	Data.Lst_asw.PushBack(answer)
}

func (Data *Data) Manage_sendball(requete *list.Element, Tab_wd *owm.All_data) {
	rqt := requete.Value.(protocol.Request)
	eball := Data.Lst_ball.Get_ballbyid(rqt.Spec.(protocol.Send_ball).Id)
	var checkpoint ballon.Checkpoints
	var answer []byte

	if eball != nil {
		eball.Value.(*ballon.Ball).Possessed = nil
		checkpoint.Coord.Longitude = rqt.Spec.(protocol.Send_ball).Lonuser
		checkpoint.Coord.Latitude = rqt.Spec.(protocol.Send_ball).Latuser
		eball.Value.(*ballon.Ball).Coord = eball.Value.(*ballon.Ball).Checkpoints.PushBack(checkpoint)
		eball.Value.(*ballon.Ball).Get_checkpointList(Tab_wd.Get_Paris())
		answer = Manage_ack(rqt.Rtype, rqt.Deviceid, eball.Value.(*ballon.Ball).Id_ball, int32(1))
		fmt.Println("Coucou :D parfait d'aller la LALALALALALAL")
	} else {
		fmt.Println("Coucou :D c'est pas bien d'aller la LALALALALALAL")
		answer = Manage_ack(rqt.Rtype, rqt.Deviceid, rqt.Spec.(protocol.Send_ball).Id, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
}

/*
** Get_answer fournis une reponse approprie a la requete du client,
** avec un buffer de 1024 Octets. Elle initialisera l'authentification de
** de l'utilisateur et nettoiera le flux de requetes traites.
 */
func (Data *Data) Get_answer(Tab_wd *owm.All_data) (er error) {
	request := Data.Lst_req.Front()
	er = nil
	if request == nil {
		er = errors.New("Get answer, but no request.")
	} else {
		Data.User, er = Data.Lst_users.Check_user(request)
		if er == nil {
			switch request.Value.(protocol.Request).Rtype {
			case SYNC:
				Data.Manage_sync(request)
			case UPDATE:
				Data.Manage_update(request)
			case POS:
				Data.Manage_pos(request)
			case TAKEN:
				Data.Manage_taken(request)
			case FOLLOW_ON:
				Data.Manage_followon(request)
			case FOLLOW_OFF:
				Data.Manage_followoff(request)
			case NEW_BALL:
				Data.Manage_newball(request, Tab_wd)
			case SEND_BALL:
				Data.Manage_sendball(request, Tab_wd)
			case ACK:
			}
		}
	}
	Del_request_done(Data.Lst_req)
	return er
}

/*
** Create a new
 */
func (Data *Data) Get_aknowledgement(Lst_usr *users.All_users) (answer []byte) {
	elem := Data.Lst_req.Back()
	treq := elem.Value.(protocol.Request)

	if treq.Rtype == NEW_BALL {
		answer = Manage_ack(treq.Rtype, treq.Deviceid, 0, int32(1))
	} else {
		answer = Manage_ack(treq.Rtype, treq.Deviceid, treq.Spec.(protocol.Ballid).Id, int32(1))
	}
	return answer
}
