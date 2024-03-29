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

import (
	"Wibo/ballon"
	"Wibo/ballonwork"
	"Wibo/devices"
	"Wibo/owm"
	"Wibo/protocol"
	"Wibo/users"
	"bytes"
	"container/list"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	_ = iota
	// SAME DEFINE FOR CLIENT AND SERVER
	ACK   = 32767
	TAKEN = 4
	// DEFINE CLIENT
	SYNC          = 1
	UPDATE        = 2
	POS           = 3
	FOLLOW_ON     = 5
	FOLLOW_OFF    = 6
	NEW_BALL      = 7
	SEND_BALL     = 8
	MAGNET        = 9
	WORKBALL      = 10
	TYPELOG       = 11 // Account connexion
	CREATEACCOUNT = 12 // Account creation without mail confirm
	SYNCROACCOUNT = 13
	DELOG         = 14 // Account diconnect
	STATSUSER     = 15
	STATSBALL     = 16
	// DEFINE SERVER
	CONN     = 1
	INF_BALL = 2
	NEARBY   = 3
	// DEFINE STATUS LOG
	UNKNOWN     = 1
	DEFAULTUSER = 2
	USERLOGGED  = 3
	// Size Packet
	SIZE_PACKET        = 1024
	SIZE_HEADER        = 16
	SIZE_STATUSER      = 104
	SIZE_STATBALL      = 56
	SIZE_COORDSTATBALL = 24
	SIZE_CONTENTBALL   = 48
)

type Posball struct {
	id       int64
	title    string
	FlagPoss int16
	lon      float64
	lat      float64
	wins     float64 //winSpeed
	wind     float64 //winDegress
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
	mess      string
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
	Lst_req     *list.List           /* Value: (*protocol.Request) which defines list request */
	Lst_asw     *list.List           /* Value: ([]byte) which defines list answer */
	Lst_ball    *ballon.All_ball     /* Value: *ballon.Ball */
	Lst_users   *users.All_users     /* Value: *users.User */
	Lst_devices *devices.All_Devices /* Value: *devices.Device */
	Lst_work    *ballonwork.All_work /* Value: *ballonwork.WorkBall */
	Logged      int16                /* Define status connexion */
	Device      *list.Element        /* Value: (*list.Element).Value.(*device.Device) */
	User        *list.Element        /* Value: (*list.Element).Value.(*users.User) */
	Logger      *log.Logger
	Conn        net.Conn
}

/* Cut one string in two news strings */
func Str_cut_n(str string, n int) (s1 string, s2 string) {
	Bstr := bytes.NewBufferString(str)
	Bs1 := Bstr.Next(n)
	s2 = Bstr.String()
	tmp := bytes.NewBuffer(Bs1)
	s1 = tmp.String()
	return s1, s2
}

/* Remove ball from followed list of User */
func RemoveBallFollowed(eball *list.Element, usr *list.Element) bool {
	user := usr.Value.(*users.User)
	for e := user.Followed.Front(); e != nil; e = e.Next() {
		if e.Value.(*list.Element) == eball {
			user.Followed.Remove(e)
			return true
		}
	}
	return false
}

/* Remove ball from possessed list of User */
func RemoveBallPossessed(eball *list.Element, usr *list.Element) bool {
	user := usr.Value.(*users.User)
	for e := user.Possessed.Front(); e != nil; e = e.Next() {
		if e.Value.(*list.Element) == eball {
			user.Possessed.Remove(e)
			return true
		}
	}
	return false
}

/* Remove user from followes list of ball */
func RemoveUserFollower(usr *list.Element, eball *list.Element) bool {
	ball := eball.Value.(*ballon.Ball)
	for e := ball.Followers.Front(); e != nil; e = e.Next() {
		if e.Value.(*list.Element) == usr {
			ball.Followers.Remove(e)
			return true
		}
	}
	return false
}

/* Checked if request list is completed */
func (Data *Data) Check_lstrequest() bool {
	Req := Data.Lst_req.Front()
	next := Req.Next()
	tmp := Req
	nr := new(protocol.Request)
	tr := new(protocol.Request)

	if Data.Logged == UNKNOWN {
		return true
	}
	tr = tmp.Value.(*protocol.Request)
	for next != nil {
		nr = next.Value.(*protocol.Request)
		if tr.Nbrpck == tr.Numpck+1 {
			return true
		} else if tr.Rtype != nr.Rtype {
			return false
		} else if tr.Numpck != nr.Numpck+1 {
			return false
		}
		next = next.Next()
		tmp = tmp.Next()
		tr = tmp.Value.(*protocol.Request)
	}
	if tr.Nbrpck == tr.Numpck+1 {
		return true
	}
	return false
}

func Del_request_done(Lst_req *list.List) {
	elem := Lst_req.Front()
	for elem != nil {
		if elem.Value.(*protocol.Request).Numpck == elem.Value.(*protocol.Request).Nbrpck-1 {
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
		binary.Write(Buffer, binary.BigEndian, pck.ptype.(*Conn).nbrball)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
		iball := pck.ptype.(*Conn).infoballs.Front()
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
func (Data *Data) Manage_ack(Type int16, IdBallon int64, value int32) (answer []byte) {
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

func Write_nearby(Req *list.Element, list_tmp *list.List, Type int16, User *users.User) (buf []byte) {
	var answer Packet
	var typesp Nearby

	answer.head.octets = 24
	answer.head.rtype = Type
	answer.head.pnbr = 1
	answer.head.pnum = 0
	typesp.nbrball = (int32)(list_tmp.Len())
	typesp.balls = list_tmp
	answer.ptype = typesp
	elem := list_tmp.Front()
	for elem != nil {
		answer.head.octets += 64
		elem = elem.Next()
	}
	Buffer := Write_header(answer)
	binary.Write(Buffer, binary.BigEndian, answer.ptype.(Nearby).nbrball)
	if User.MagnetisValid() == true {
		binary.Write(Buffer, binary.BigEndian, int32(1))
	} else {
		binary.Write(Buffer, binary.BigEndian, int32(0))
	}
	elem = answer.ptype.(Nearby).balls.Front()
	for elem != nil {
		ifb := elem.Value.(Posball)
		binary.Write(Buffer, binary.BigEndian, ifb.id)
		Buffer.WriteString(ifb.title)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 16-len(ifb.title)))
		binary.Write(Buffer, binary.BigEndian, ifb.FlagPoss)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 6))
		binary.Write(Buffer, binary.BigEndian, ifb.lon)
		binary.Write(Buffer, binary.BigEndian, ifb.lat)
		binary.Write(Buffer, binary.BigEndian, ifb.wins)
		binary.Write(Buffer, binary.BigEndian, ifb.wind)
		elem = elem.Next()
	}
	binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-answer.head.octets))
	buf = Buffer.Bytes()
	return buf
}

func GetPacketsContent(ball *ballon.Ball, typeR int16) *list.List {
	lstP := list.New()
	pck := new(Packet)

	pck.head.octets = SIZE_HEADER + SIZE_CONTENTBALL
	pck.head.rtype = typeR
	pck.head.pnum = 0
	pck.ptype = new(Contentball)
	pck.ptype.(*Contentball).messages = list.New()
	for em := ball.Messages.Front(); em != nil; em = em.Next() {
		msg := em.Value.(ballon.Message)
		if int16(msg.Size)+16+pck.head.octets > SIZE_PACKET {
			pck.ptype.(*Contentball).nbruser = int32(ball.Stats.NbrFollow)
			pck.ptype.(*Contentball).nbrmess = int32(pck.ptype.(*Contentball).messages.Len())
			fmt.Println("Number message of packet: ", pck.ptype.(*Contentball).nbrmess)
			lstP.PushBack(pck)
			tmp := pck
			pck = new(Packet)
			pck.head.octets = SIZE_HEADER + SIZE_CONTENTBALL
			pck.head.rtype = typeR
			pck.head.pnum = tmp.head.pnum + 1
			pck.ptype = new(Contentball)
			pck.ptype.(*Contentball).messages = list.New()
		}
		var mes Message
		mes.size = msg.Size
		mes.id = msg.Id
		fmt.Println("MessageIdPack:", mes.id)
		fmt.Println("MessageIdSource:", msg.Id)
		mes.idcountry = msg.Idcountry
		mes.idcity = msg.Idcity
		mes.mess = msg.Content
		mes.mtype = msg.Type
		pck.head.octets += int16(16 + msg.Size)
		pck.ptype.(*Contentball).messages.PushBack(mes)
	}
	pck.ptype.(*Contentball).nbruser = int32(ball.Stats.NbrFollow)
	pck.ptype.(*Contentball).nbrmess = int32(pck.ptype.(*Contentball).messages.Len())
	fmt.Println("Number message of packet: ", pck.ptype.(*Contentball).nbrmess)
	lstP.PushBack(pck)
	Pnbr := lstP.Len()
	for e := lstP.Front(); e != nil; e = e.Next() {
		e.Value.(*Packet).head.pnbr = int32(Pnbr)
	}
	return lstP
}

func Write_contentball(Ball *ballon.Ball, packettype int16) (Alst *list.List) {
	Alst = list.New()
	lstPack := GetPacketsContent(Ball, packettype)

	year := int16(Ball.Stats.CreationDate.Year())
	month := int16(Ball.Stats.CreationDate.Month())
	day := int16(Ball.Stats.CreationDate.Day())
	houre := int16(Ball.Stats.CreationDate.Hour())
	minute := int16(Ball.Stats.CreationDate.Minute())
	sizeTitle := len(Ball.Title)
	fmt.Println("Title + size: ", Ball.Title)
	fmt.Println(sizeTitle)

	for ep := lstPack.Front(); ep != nil; ep = ep.Next() {
		pck := ep.Value.(*Packet)
		Buffer := Write_header(*pck)
		binary.Write(Buffer, binary.BigEndian, Ball.Id_ball)
		binary.Write(Buffer, binary.BigEndian, year)
		binary.Write(Buffer, binary.BigEndian, month)
		binary.Write(Buffer, binary.BigEndian, day)
		binary.Write(Buffer, binary.BigEndian, houre)
		binary.Write(Buffer, binary.BigEndian, minute)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 6))
		Buffer.WriteString(Ball.Title)
		fmt.Println("Title ball: ", Ball.Title)
		if 16-sizeTitle > 0 {
			binary.Write(Buffer, binary.BigEndian, make([]byte, 16-sizeTitle))
		}
		binary.Write(Buffer, binary.BigEndian, pck.ptype.(*Contentball).nbruser)
		binary.Write(Buffer, binary.BigEndian, pck.ptype.(*Contentball).nbrmess)
		for em := pck.ptype.(*Contentball).messages.Front(); em != nil; em = em.Next() {
			mes := em.Value.(Message)
			binary.Write(Buffer, binary.BigEndian, mes.id)
			binary.Write(Buffer, binary.BigEndian, mes.size)
			binary.Write(Buffer, binary.BigEndian, mes.idcountry)
			binary.Write(Buffer, binary.BigEndian, mes.idcity)
			Buffer.WriteString(mes.mess)
		}
		binary.Write(Buffer, binary.BigEndian, make([]byte, SIZE_PACKET-pck.head.octets))
		Alst.PushBack(Buffer.Bytes())
	}
	return Alst
}

func (Data *Data) Write_StatUser(year, month, day, houre, minute int16) (buf []byte) {
	user := Data.User.Value.(*users.User)
	var answer Packet

	answer.head.octets = SIZE_HEADER + SIZE_STATUSER
	answer.head.rtype = STATSUSER
	answer.head.pnbr = 1
	answer.head.pnum = 0
	Buffer := Write_header(answer)
	binary.Write(Buffer, binary.BigEndian, year)
	binary.Write(Buffer, binary.BigEndian, month)
	binary.Write(Buffer, binary.BigEndian, day)
	binary.Write(Buffer, binary.BigEndian, houre)
	binary.Write(Buffer, binary.BigEndian, minute)
	binary.Write(Buffer, binary.BigEndian, make([]byte, 6))
	binary.Write(Buffer, binary.BigEndian, user.Stats.NbrBallCreate)
	binary.Write(Buffer, binary.BigEndian, user.Stats.NbrCatch)
	binary.Write(Buffer, binary.BigEndian, user.Stats.NbrSend)
	binary.Write(Buffer, binary.BigEndian, user.Stats.NbrFollow)
	binary.Write(Buffer, binary.BigEndian, user.Stats.NbrMessage)
	binary.Write(Buffer, binary.BigEndian, Data.Lst_users.NbrUsers)
	binary.Write(Buffer, binary.BigEndian, Data.Lst_users.GlobalStat.NbrBallCreate)
	binary.Write(Buffer, binary.BigEndian, Data.Lst_users.GlobalStat.NbrCatch)
	binary.Write(Buffer, binary.BigEndian, Data.Lst_users.GlobalStat.NbrSend)
	binary.Write(Buffer, binary.BigEndian, Data.Lst_users.GlobalStat.NbrFollow)
	binary.Write(Buffer, binary.BigEndian, Data.Lst_users.GlobalStat.NbrMessage)
	binary.Write(Buffer, binary.BigEndian, make([]byte, 1024-answer.head.octets))
	buf = Buffer.Bytes()
	return buf
}

/* Get all follower ball and create packets list for send the infoballs */
func (Data *Data) Manage_sync(Req *list.Element) {
	flwupball := Data.User.Value.(*users.User).Followed.Front()
	plist := list.New()
	pck := new(Packet)
	conn := new(Conn)

	pck.head.octets = SIZE_HEADER + 8
	pck.head.rtype = CONN
	conn.nbrball = 0
	conn.infoballs = list.New()

	for flwupball != nil {
		myball := new(Infoball)
		ball := flwupball.Value.(*list.Element).Value.(*ballon.Ball)
		myball.id = ball.Id_ball
		myball.nbrmess = int32(ball.Messages.Len())
		if ball.Possessed == Data.User {
			myball.taken = 1
		} else {
			myball.taken = 0
		}
		if pck.head.octets+16 > SIZE_PACKET {
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
	pck.ptype = conn
	plist.PushBack(pck)
	pelem := plist.Front()
	nbrpacket := int32(plist.Len())
	numpacket := int32(0)
	for pelem != nil {
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

	rqt := request.Value.(*protocol.Request)
	idsearch := rqt.Spec.(protocol.Ballid).Id
	eball := Data.User.Value.(*users.User).Followed.Front()
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
		answer := Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Ballid).Id, int32(0))
		Data.Lst_asw.PushBack(answer)
	}
}

func (Data *Data) Manage_pos(Req *list.Element) {
	list_tmp := list.New()
	var ifball Posball

	eball := Data.Lst_ball.Blist.Front()
	for eball != nil {
		ball := eball.Value.(*ballon.Ball)
		if ball.Check_nearbycoord(Req) == true {
			ball := eball.Value.(*ballon.Ball)
			ifball.id = ball.Id_ball
			ifball.title = ball.Title
			if ball.Check_userCreated(Data.User) == true {
				ifball.FlagPoss = 1
			} else if ball.Check_userfollower(Data.User) == true {
				ifball.FlagPoss = 2
			} else {
				ifball.FlagPoss = 0
			}
			ifball.lon = ball.Coord.Value.(ballon.Checkpoint).Coord.Lon
			ifball.lat = ball.Coord.Value.(ballon.Checkpoint).Coord.Lat
			ifball.wins = ball.Wind.Speed
			ifball.wind = ball.Wind.Degress
			list_tmp.PushBack(ifball)
		}
		eball = eball.Next()
	}
	Len := list_tmp.Len()
	for Len > 10 {
		elem := list_tmp.Front()
		random := rand.Intn(Len)
		for elem != nil && random > 0 {
			elem = elem.Next()
			random -= 1
		}
		list_tmp.Remove(eball)
		Len -= 1
	}
	answer := Write_nearby(Req, list_tmp, POS, Data.User.Value.(*users.User))
	Data.Lst_asw.PushBack(answer)
}

func (Data *Data) Manage_taken(request *list.Element, Wd *owm.All_data) {
	eball := Data.Lst_ball.Blist.Front()
	rqt := request.Value.(*protocol.Request)
	user := Data.User.Value.(*users.User)

	for eball != nil && eball.Value.(*ballon.Ball).Id_ball != request.Value.(*protocol.Request).Spec.(protocol.Taken).Id {
		eball = eball.Next()
	}
	if eball != nil && user.Possessed.Len() < 3 {
		ball := eball.Value.(*ballon.Ball)
		if ball.Possessed == nil && ball.Check_userfollower(Data.User) == false && ball.Check_userCreated(Data.User) == false {
			ball.Possessed = Data.User
			ball.Edited = true
			ball.Followers.PushBack(Data.User)
			user.Followed.PushBack(eball)
			user.Possessed.PushBack(eball)
			/* Begin Stats */
			user.Stats.NbrCatch++
			user.Stats.NbrFollow++
			Data.Lst_users.GlobalStat.NbrCatch++
			if ball.Followers.Len() == 1 {
				Data.Lst_users.GlobalStat.NbrFollow++
			}
			ball.Stats.NbrCatch++
			ball.Stats.NbrFollow++
			ball.Stats.NbrKm += ball.GetDistance(rqt.Coord.Lon, rqt.Coord.Lat)
			ball.InitCoord(rqt.Coord.Lon, rqt.Coord.Lat, rqt.Spec.(protocol.Taken).FlagMagnet, Wd, false)
			if rqt.Spec.(protocol.Taken).FlagMagnet == 1 {
				Data.User.Value.(*users.User).Magnet = time.Now()
				ball.Stats.NbrMagnet++
			}
			Lst_answer := Write_contentball(ball, TAKEN)
			Data.Lst_asw.PushBackList(Lst_answer)
			/* End Stats */
		} else {
			answer := Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Taken).Id, int32(0))
			Data.Lst_asw.PushBack(answer)
		}
	} else {
		answer := Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Taken).Id, int32(0))
		Data.Lst_asw.PushBack(answer)
	}
}

func (Data *Data) Manage_followon(request *list.Element) {
	rqt := request.Value.(*protocol.Request)

	eball := Data.Lst_ball.Blist.Front()
	for eball != nil && eball.Value.(*ballon.Ball).Id_ball != rqt.Spec.(protocol.Ballid).Id {
		eball = eball.Next()
	}
	var answer []byte
	if eball != nil && eball.Value.(*ballon.Ball).Check_userfollower(Data.User) == false {
		eball.Value.(*ballon.Ball).Edited = true
		eball.Value.(*ballon.Ball).Followers.PushBack(Data.User)
		user := Data.User.Value.(*users.User)
		user.Followed.PushBack(eball)
		user.Stats.NbrFollow++
		if eball.Value.(*ballon.Ball).Followers.Len() == 1 {
			Data.Lst_users.GlobalStat.NbrFollow++
		}
		eball.Value.(*ballon.Ball).Stats.NbrFollow++
		answer = Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Ballid).Id, int32(1))
	} else {
		answer = Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Ballid).Id, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
}

func (Data *Data) Manage_followoff(request *list.Element) {
	var answer []byte
	rqt := request.Value.(*protocol.Request)
	eball := Data.Lst_ball.Blist.Front()

	for eball != nil &&
		eball.Value.(*ballon.Ball).Id_ball != rqt.Spec.(protocol.Ballid).Id {
		eball = eball.Next()
	}
	if eball != nil &&
		eball.Value.(*ballon.Ball).Check_userfollower(Data.User) == true {
		RemoveUserFollower(Data.User, eball)
		eball.Value.(*ballon.Ball).Edited = true
		user := Data.User.Value.(*users.User)
		RemoveBallFollowed(eball, Data.User)
		user.Stats.NbrFollow--
		if eball.Value.(*ballon.Ball).Followers.Len() == 0 {
			Data.Lst_users.GlobalStat.NbrFollow--
		}
		eball.Value.(*ballon.Ball).Stats.NbrFollow--
		answer = Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Ballid).Id, int32(1))

	} else {
		answer = Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Ballid).Id, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
}

func (Data *Data) Manage_newball(requete *list.Element, Tab_wd *owm.All_data) {
	ball := new(ballon.Ball)
	ball.Stats = new(ballon.StatsBall)
	rqt := requete.Value.(*protocol.Request)
	var newball protocol.New_ball
	var mess ballon.Message
	user := Data.User.Value.(*users.User)

	if user.NbrBallSend < 10 {
		newball = requete.Value.(*protocol.Request).Spec.(protocol.New_ball)
		ball.Id_ball = Data.Lst_ball.Id_max
		Data.Lst_ball.Id_max++
		ball.FlagC = true
		ball.Edited = true
		ball.Title = newball.Title
		ball.Messages = list.New()
		ball.Itinerary = list.New()
		ball.Followers = list.New()
		ball.Checkpoints = list.New()
		mess.Id = 0
		mess.Size = (int32)(len(newball.Message))
		mess.Content = newball.Message
		mess.Type = 1
		ball.Messages.PushFront(mess)
		ball.Date = time.Now()
		ball.Possessed = nil
		ball.Followers.PushFront(Data.User)
		ball.Creator = Data.User
		ball.InitCoord(rqt.Coord.Lon, rqt.Coord.Lat, int16(0), Tab_wd, true)
		eball := Data.Lst_ball.Blist.PushBack(ball)
		user := Data.User.Value.(*users.User)
		user.Followed.PushBack(eball)
		user.NbrBallSend++
		// Begin Stats
		user.Stats.NbrMessage++
		user.Stats.NbrBallCreate++
		user.Stats.NbrSend++
		user.Stats.NbrFollow++
		Data.Lst_users.GlobalStat.NbrBallCreate++
		Data.Lst_users.GlobalStat.NbrSend++
		Data.Lst_users.GlobalStat.NbrFollow++
		Data.Lst_users.GlobalStat.NbrMessage++
		ball.Stats.CreationDate = time.Now()
		ball.Stats.CoordCreated = new(ballon.Coordinate)
		ball.Stats.CoordCreated.Lon = rqt.Coord.Lon
		ball.Stats.CoordCreated.Lat = rqt.Coord.Lat
		ball.Stats.NbrFollow++
		// End Stats
		answer := Data.Manage_ack(rqt.Rtype, ball.Id_ball, int32(1))
		Data.Lst_asw.PushBack(answer)
	} else {
		answer := Data.Manage_ack(rqt.Rtype, 0, int32(0))
		Data.Lst_asw.PushBack(answer)
	}
}

func (Data *Data) Manage_sendball(requete *list.Element, Tab_wd *owm.All_data) {
	rqt := requete.Value.(*protocol.Request)
	eball := Data.Lst_ball.Get_ballbyid(rqt.Spec.(protocol.Send_ball).Id)
	var answer []byte

	if eball != nil && eball.Value.(*ballon.Ball).Check_userPossessed(Data.User) == true {
		user := Data.User.Value.(*users.User)
		ball := eball.Value.(*ballon.Ball)
		RemoveBallPossessed(eball, Data.User)
		ball.Possessed = nil
		ball.Edited = true
		ball.InitCoord(rqt.Coord.Lon, rqt.Coord.Lat, int16(0), Tab_wd, true)
		var message ballon.Message
		message.Id = int32(ball.Messages.Len())
		message.Size = rqt.Spec.(protocol.Send_ball).Octets
		message.Content = rqt.Spec.(protocol.Send_ball).Message
		message.Type = 1
		ball.Messages.PushBack(message)
		/* Begin stats ---- */
		user.Stats.NbrMessage++
		user.Stats.NbrSend++
		Data.Lst_users.GlobalStat.NbrSend++
		Data.Lst_users.GlobalStat.NbrMessage++
		/* End Stats --- */
		answer = Data.Manage_ack(rqt.Rtype, eball.Value.(*ballon.Ball).Id_ball, int32(1))
	} else {
		answer = Data.Manage_ack(rqt.Rtype, rqt.Spec.(protocol.Send_ball).Id, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
}

/* Create tree ball list with id random. If Ball is already taked, a first ball next is taked */
func (Data *Data) Manage_magnet(requete *list.Element, Tab_wd *owm.All_data) {
	rqt := requete.Value.(*protocol.Request)
	var tab [3]int64
	list_tmp := list.New()
	var ifball Posball
	user := Data.User.Value.(*users.User)
	var answer []byte
	var eball *list.Element

	if user.MagnetisValid() == true {
		if Data.Lst_ball.Id_max > 0 {
			for i := 0; i < 3; i++ {
				tab[i] = rand.Int63n(Data.Lst_ball.Id_max)
			}
			list_tmp_2 := Data.Lst_ball.Get_ballbyid_tomagnet(tab, Data.User)
			eball = list_tmp_2.Front()
		}
		for eball != nil {
			ball := eball.Value.(*list.Element).Value.(*ballon.Ball)
			ifball.id = ball.Id_ball
			ifball.title = ball.Title
			ifball.FlagPoss = 0
			ifball.lon = ball.Coord.Value.(ballon.Checkpoint).Coord.Lon
			ifball.lat = ball.Coord.Value.(ballon.Checkpoint).Coord.Lat
			ifball.wins = ball.Wind.Speed
			ifball.wind = ball.Wind.Degress
			list_tmp.PushBack(ifball)
			eball = eball.Next()
		}
		answer = Write_nearby(requete, list_tmp, MAGNET, user)
	} else {
		answer = Data.Manage_ack(rqt.Rtype, 0, 0)
	}
	Data.Lst_asw.PushBack(answer)
}

func (Data *Data) Manage_Login(request *list.Element, Db *sql.DB, Dlist *devices.All_Devices) (er error) {
	req := request.Value.(*protocol.Request)
	er = nil
	flag := true
	var answer []byte

	if req.Rtype != TYPELOG {
		er = errors.New("Bad type to Manage_Login")
		answer = Data.Manage_ack(TYPELOG, 0, int32(0))
	} else {
		if Data.Device == nil {
			Data.Device, er = Dlist.GetDevice(request, Db, Data.Lst_users)
		}
		if er == nil {
			device := Data.Device.Value.(*devices.Device)
			if len(req.Spec.(protocol.Log).Email) <= 1 {
				Data.Logged = DEFAULTUSER
				Data.User = device.UserDefault
			} else {
				Data.User = Data.Lst_users.Check_user(request, Db, device.Historic)
				if Data.User == nil {
					flag = false
					Data.Logged = DEFAULTUSER
					device.UserSpec = nil
					Data.User = device.UserDefault
				} else {
					Data.Logged = USERLOGGED
					device.UserSpec = Data.User
					device.AddUserSpecOnHistory(Data.User)
				}
			}
		}
		if er != nil {
			answer = Data.Manage_ack(TYPELOG, 0, int32(0))
		} else if flag == false {
			answer = Data.Manage_ack(TYPELOG, 0, int32(2))
		} else {
			answer = Data.Manage_ack(TYPELOG, 0, int32(1))
		}
	}
	Data.Lst_asw.PushBack(answer)
	return er
}

func (Data *Data) Manage_CreateAccount(request *list.Element, Db *sql.DB) (er error) {
	req := request.Value.(*protocol.Request)
	er = nil
	User := new(users.User)
	var answer []byte

	if users.CheckValidMail(req.Spec.(protocol.Log).Email) == true {
		User.Mail = req.Spec.(protocol.Log).Email
		User.NbrBallSend = 0
		User.Coord.Lon = req.Coord.Lon
		User.Log = time.Now()
		User.Coord.Lat = req.Coord.Lat
		User.Followed = list.New()
		User.Possessed = list.New()
		User.HistoricReq = list.New()
		User.Stats = new(users.StatsUser)
		User.Stats.CreationDate = time.Now()
		User.Psd = req.Spec.(protocol.Log).Pswd
		flag, err := Data.Lst_users.Add_new_user(User, Db, req.Spec.(protocol.Log).Pswd)
		er = err
		if err != nil {
			Data.Logger.Println("Error Add_new_user: ", err)
		}
		if flag == true {
			eUser := Data.Lst_users.Ulist.PushFront(User)
			Data.Device.Value.(*devices.Device).Historic.PushFront(eUser)
			/* Begin Stats */
			User.Stats = new(users.StatsUser)
			User.Stats.CreationDate = time.Now()
			Data.Lst_users.NbrUsers++
			/* End Stats */
			answer = Data.Manage_ack(CREATEACCOUNT, 0, int32(1))
		} else {
			er = nil
			answer = Data.Manage_ack(CREATEACCOUNT, 0, int32(0))
		}
	} else {
		answer = Data.Manage_ack(CREATEACCOUNT, 0, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
	return er
}

func AddFollowed(euser *list.Element, euserDefault *list.Element) {
	user := euser.Value.(*users.User)
	userDefault := euserDefault.Value.(*users.User)
	var tball *list.Element

	for eball := userDefault.Followed.Front(); eball != nil; eball = eball.Next() {
		ball := eball.Value.(*list.Element).Value.(*ballon.Ball)
		for tball = user.Followed.Front(); tball != nil; tball = tball.Next() {
			idball := tball.Value.(*list.Element).Value.(*ballon.Ball).Id_ball
			if idball == ball.Id_ball {
				break
			}
		}
		if tball == nil {
			user.Followed.PushBack(eball.Value.(*list.Element))
			ball.Followers.PushFront(euser)
		}
	}
}

func GetPossessed(euser *list.Element, euserDefault *list.Element) {
	user := euser.Value.(*users.User)
	userDefault := euserDefault.Value.(*users.User)

	for eball := userDefault.Possessed.Front(); eball != nil; eball = userDefault.Possessed.Front() {
		ball := eball.Value.(*list.Element).Value.(*ballon.Ball)
		user.Possessed.PushBack(userDefault.Possessed.Remove(eball))
		ball.Possessed = euser
	}
}

func (Data *Data) Manage_SyncAccount(request *list.Element, Db *sql.DB) (er error) {
	er = nil
	req := request.Value.(*protocol.Request)
	var answer []byte

	if Data.Logged == USERLOGGED {
		device := Data.Device.Value.(*devices.Device)
		user := Data.User.Value.(*users.User)
		userDefault := device.UserDefault.Value.(*users.User)
		user.NbrBallSend += userDefault.NbrBallSend
		user.Coord.Lon = req.Coord.Lon
		user.Coord.Lat = req.Coord.Lat
		user.Log = time.Now()
		AddFollowed(device.UserSpec, device.UserDefault)
		GetPossessed(device.UserSpec, device.UserDefault)
		answer = Data.Manage_ack(SYNCROACCOUNT, 0, int32(1))
	} else {
		answer = Data.Manage_ack(SYNCROACCOUNT, 0, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
	return er
}

func (Data *Data) Manage_Delog(request *list.Element, Db *sql.DB) (er error) {
	device := Data.Device.Value.(*devices.Device)
	var answer []byte

	if Data.Logged == USERLOGGED {
		Data.User.Value.(*users.User).Log = time.Now()
		Data.User = device.UserDefault
		device.UserSpec = nil
		answer = Data.Manage_ack(DELOG, 0, int32(1))
		Data.Logged = DEFAULTUSER
	} else {
		answer = Data.Manage_ack(DELOG, 0, int32(0))
	}
	Data.Lst_asw.PushBack(answer)
	return er
}

func (Data *Data) Manage_StatUser(request *list.Element) {
	var answer []byte
	user := Data.User.Value.(*users.User)
	if user == nil {
		answer = Data.Manage_ack(STATSUSER, 0, int32(0))
	} else {
		year := int16(user.Stats.CreationDate.Year())
		month := int16(user.Stats.CreationDate.Month())
		day := int16(user.Stats.CreationDate.Day())
		houre := int16(user.Stats.CreationDate.Hour())
		minute := int16(user.Stats.CreationDate.Minute())
		answer = Data.Write_StatUser(year, month, day, houre, minute)
	}
	Data.Lst_asw.PushBack(answer)
}

func Write_StatBall(lst *list.List, nbrCheck int32, nbrPack int, ball *ballon.Ball) *list.List {
	var answer Packet
	var NbrItin int32
	lst_asw := list.New()
	eCheck := lst.Front()

	for i := nbrPack; i > 0; i-- {
		NbrItin = int32((SIZE_PACKET - (SIZE_HEADER + SIZE_STATBALL)) / SIZE_COORDSTATBALL)
		if nbrCheck > NbrItin {
			nbrCheck -= NbrItin
		} else {
			NbrItin = nbrCheck
		}
		answer.head.octets = int16(SIZE_HEADER + SIZE_STATBALL + (NbrItin * SIZE_COORDSTATBALL))
		if i == nbrPack {
			answer.head.rtype = STATSBALL
			answer.head.pnbr = int32(nbrPack)
			answer.head.pnum = 0
		} else {
			answer.head.pnum++
		}
		Buffer := Write_header(answer)
		binary.Write(Buffer, binary.BigEndian, ball.Stats.CoordCreated.Lon)
		binary.Write(Buffer, binary.BigEndian, ball.Stats.CoordCreated.Lat)
		binary.Write(Buffer, binary.BigEndian, ball.Stats.NbrKm)
		binary.Write(Buffer, binary.BigEndian, ball.Stats.NbrFollow)
		binary.Write(Buffer, binary.BigEndian, ball.Stats.NbrCatch)
		binary.Write(Buffer, binary.BigEndian, ball.Stats.NbrMagnet)
		binary.Write(Buffer, binary.BigEndian, NbrItin)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 4))
		for j := int32(1); j <= NbrItin; j++ {
			check := eCheck.Value.(*ballon.Checkpoint)
			binary.Write(Buffer, binary.BigEndian, int32(j))
			binary.Write(Buffer, binary.BigEndian, int32(check.MagnetFlag))
			binary.Write(Buffer, binary.BigEndian, check.Coord.Lon)
			binary.Write(Buffer, binary.BigEndian, check.Coord.Lat)
			eCheck = eCheck.Next()
		}
		binary.Write(Buffer, binary.BigEndian, make([]byte, SIZE_PACKET-answer.head.octets))
		buf := Buffer.Bytes()
		lst_asw.PushBack(buf)
	}
	return lst_asw
}

func (Data *Data) Manage_StatBall(request *list.Element, Db *sql.DB) {
	rqt := request.Value.(*protocol.Request)
	eball := Data.Lst_ball.Get_ballbyid(rqt.Spec.(protocol.Ballid).Id)
	ball := eball.Value.(*ballon.Ball)
	nbrCheckpoint, LstCheckpoint := ball.GetItinerary(Db, Data.Lst_ball)

	sizeStat := (SIZE_PACKET - SIZE_STATBALL - SIZE_HEADER)
	var tmpPacket float64
	tmpPacket = (float64(SIZE_COORDSTATBALL) * float64(nbrCheckpoint)) / float64(sizeStat)
	nbrPacket := int(tmpPacket)
	if float64(nbrPacket) < tmpPacket || nbrPacket == 0 {
		nbrPacket++
	}
	lst_asw := Write_StatBall(LstCheckpoint, nbrCheckpoint, nbrPacket, ball)
	if lst_asw == nil {
		answer := Data.Manage_ack(STATSBALL, 0, int32(0))
		Data.Lst_asw.PushBack(answer)
	} else {
		Data.Lst_asw.PushBackList(lst_asw)
	}
}

func Write_workball(lst_work *list.List) *list.List {
	lst_asw := list.New()
	var answer Packet

	answer.head.octets = SIZE_PACKET
	answer.head.rtype = WORKBALL
	answer.head.pnbr = int32(lst_work.Len())
	answer.head.pnum = 0
	for e := lst_work.Front(); e != nil; e = e.Next() {
		workball := e.Value.(*ballonwork.WorkBall)
		Buffer := Write_header(answer)
		answer.head.pnum++
		Buffer.WriteString(workball.Title)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 16-len(workball.Title)))
		Buffer.WriteString(workball.Message)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 656-len(workball.Message)))
		binary.Write(Buffer, binary.BigEndian, workball.Coord.Lon)
		binary.Write(Buffer, binary.BigEndian, workball.Coord.Lat)
		Buffer.WriteString(workball.Link)
		binary.Write(Buffer, binary.BigEndian, make([]byte, 320-len(workball.Link)))
		lst_asw.PushBack(Buffer.Bytes())
	}
	return lst_asw
}

func (Data *Data) Manage_WorkBall(request *list.Element) {
	lst_work := list.New()
	for e := Data.Lst_work.Wlist.Front(); e != nil; e = e.Next() {
		workBall := e.Value.(*ballonwork.WorkBall)
		if workBall.Check_nearbycoord(request) == true {
			lst_work.PushBack(workBall)
		}
	}
	if lst_work.Len() == 0 {
		answer := Data.Manage_ack(WORKBALL, 0, int32(0))
		Data.Lst_asw.PushBack(answer)
	} else {
		Data.Lst_asw.PushBackList(Write_workball(lst_work))
	}
}

func (Data *Data) Get_answer(Tab_wd *owm.All_data, Db *sql.DB) (er error) {
	request := Data.Lst_req.Front()
	er = nil

	if Data.Logged == UNKNOWN {
		er = Data.Manage_Login(request, Db, Data.Lst_devices)
	} else {
		switch request.Value.(*protocol.Request).Rtype {
		default:
			er = errors.New("Get_answer detect invalid type")
		case SYNC:
			Data.Manage_sync(request)
		case UPDATE:
			Data.Manage_update(request)
		case POS:
			Data.Manage_pos(request)
		case TAKEN:
			Data.Manage_taken(request, Tab_wd)
		case FOLLOW_ON:
			Data.Manage_followon(request)
		case FOLLOW_OFF:
			Data.Manage_followoff(request)
		case NEW_BALL:
			Data.Manage_newball(request, Tab_wd)
		case SEND_BALL:
			Data.Manage_sendball(request, Tab_wd)
		case MAGNET:
			Data.Manage_magnet(request, Tab_wd)
		case WORKBALL:
			Data.Manage_WorkBall(request)
		case ACK:
		case TYPELOG:
			er = Data.Manage_Login(request, Db, Data.Lst_devices)
		case CREATEACCOUNT:
			er = Data.Manage_CreateAccount(request, Db)
		case SYNCROACCOUNT:
			er = Data.Manage_SyncAccount(request, Db)
		case DELOG:
			er = Data.Manage_Delog(request, Db)
		case STATSUSER:
			Data.Manage_StatUser(request)
		case STATSBALL:
			Data.Manage_StatBall(request, Db)
		}
	}
	Del_request_done(Data.Lst_req)
	return er
}

/*
** To perform exchange's multiple. Call from sock.go "handleConnection"
 */
func (Data *Data) Get_aknowledgement(Lst_usr *users.All_users) (answer []byte) {
	elem := Data.Lst_req.Back()
	treq := elem.Value.(*protocol.Request)

	if treq.Rtype == NEW_BALL {
		answer = Data.Manage_ack(treq.Rtype, 0, int32(1))
	} else {
		answer = Data.Manage_ack(treq.Rtype, treq.Spec.(protocol.Ballid).Id, int32(1))
	}
	return answer
}
