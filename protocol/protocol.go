//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  protocol.go                                        :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  By: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  Created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  Updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	_             = iota
	ACK           = 32767
	SYNC          = 1
	MAJ           = 2
	POS           = 3
	TAKEN         = 4
	FOLLOW_ON     = 5
	FOLLOW_OFF    = 6
	NEW_BALL      = 7
	SEND_BALL     = 8
	MAGNET        = 9
	ITINERARY     = 10
	TYPELOG       = 11
	CREATEACCOUNT = 12
	SYNCROACCOUNT = 13
	DELOG         = 14
)

type Ack struct {
	Atype  int16
	Status int16
}

type Position struct {
	Lon float64
	Lat float64
}

type Ballid struct {
	Id int64
}

type New_ball struct {
	Title   string
	Octets  int32
	Message string
}

type Send_ball struct {
	Id      int64
	Octets  int32
	Message string
}

type Log struct {
	Email string
	Pswd  string
}

//iddevice int64 // Deviendra une string ou un buffer ..
type Request struct {
	Octets   int16
	Rtype    int16
	Nbrpck   int32
	Numpck   int32
	IdMobile string
	//	Deviceid int64
	Coord Position
	Spec  interface{}
}

/* Decode type Ack */
func Request_ack(TypBuff *bytes.Buffer) (ack Ack, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &ack.Atype)
	if err != nil {
		er = errors.New("Get_ack in protocol: Error binary.Read")
		return ack, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &ack.Status)
	if err != nil {
		er = errors.New("Get_ack in protocol: Error binary.Read")
		return ack, er
	}
	return ack, er
}

/* Decode types MAJ, TAKEN, FOLLOW_ON, and FOLLOW_OFF */
func Request_idball(TypBuff *bytes.Buffer) (id Ballid, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &id.Id)

	if err != nil {
		er = errors.New("Get_idball in protocol: Error binary.Read")
		return id, er
	}
	return id, er
}

/* Decode type newball */
func Request_newball(TypBuff *bytes.Buffer) (ball New_ball, er error) {
	var err error

	ball.Title, err = TypBuff.ReadString(0)
	if len(ball.Title) == 1 {
		er = errors.New("Get_newball in protocol 1: Error ReadString")
		return ball, er
	}
	TypBuff.Next((16 - len(ball.Title)))
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Octets)
	if err != nil {
		er = errors.New("Get_newball in protocol 2: Error binary.Read")
		return ball, er
	}
	TypBuff.Next(4)
	ball.Message, err = TypBuff.ReadString(0)
	if 1 == len(ball.Message) {
		er = errors.New("Get_newball in protocol 3: Error ReadString")
		return ball, er
	}
	return ball, er
}

/* Decode type sendball  */
func Request_sendball(TypBuff *bytes.Buffer) (ball Send_ball, er error) {
	var err error

	err = binary.Read(TypBuff, binary.BigEndian, &ball.Id)
	if err != nil {
		er = errors.New("Get_sendball in protocol: Error binary.Read")
		return ball, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Octets)
	if err != nil {
		er = errors.New("Get_sendball in protocol: Error binary.Read")
		return ball, er
	}
	TypBuff.Next(4)
	ball.Message, err = TypBuff.ReadString(0)
	if 1 == len(ball.Message) {
		er = errors.New("Get_sendball in protocol: Error ReadString")
		return ball, er
	}
	return ball, er
}

func Request_Log(TypBuff *bytes.Buffer) (log Log, er error) {
	var err error
	var email [320]byte
	var pswd [512]byte

	err = binary.Read(TypBuff, binary.BigEndian, &email)
	if err != nil {
		er = errors.New("Get_sendball in protocol: Error binary.Read")
		return log, er
	}
	log.Email = string(email[:320])
	err = binary.Read(TypBuff, binary.BigEndian, &pswd)
	if err != nil {
		er = errors.New("Get_sendball in protocol: Error binary.Read")
		return log, er
	}
	log.Pswd = string(pswd[:512])
	return log, er
}

/*
** Master function that tokenizes client's request.
** See the protocol
 */
func (token *Request) Get_request(buff []byte) (er error) {
	TypBuff := bytes.NewBuffer(buff)
	var IdMobile [40]byte

	er = nil
	err := binary.Read(TypBuff, binary.BigEndian, &token.Octets)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &token.Rtype)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	TypBuff.Next(4)
	err = binary.Read(TypBuff, binary.BigEndian, &token.Nbrpck)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &token.Numpck)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &IdMobile)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	token.IdMobile = string(IdMobile[:40])
	//	TypBuff.Next(32)
	err = binary.Read(TypBuff, binary.BigEndian, &token.Coord.Lon)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &token.Coord.Lat)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return er
	}
	switch token.Rtype {
	case SYNC:
		return er
	case MAJ, TAKEN, FOLLOW_ON, FOLLOW_OFF, ITINERARY:
		token.Spec, er = Request_idball(TypBuff)
	case POS:
	case NEW_BALL:
		token.Spec, er = Request_newball(TypBuff)
	case SEND_BALL:
		token.Spec, er = Request_sendball(TypBuff)
	case ACK:
		token.Spec, er = Request_ack(TypBuff)
	case TYPELOG, CREATEACCOUNT:
		token.Spec, er = Request_Log(TypBuff)
	case SYNCROACCOUNT:
	case DELOG:
	}
	return er
}

/* Print request to debug */
func (token *Request) Print_token_debug() {
	fmt.Println("Request header:")
	fmt.Println(token.Octets)
	fmt.Println(token.Rtype)
	fmt.Println(token.Nbrpck)
	fmt.Println(token.Numpck)
	fmt.Println(token.IdMobile)
	fmt.Println(token.Coord.Lon)
	fmt.Println(token.Coord.Lat)
	fmt.Println("Type request:")
	switch token.Rtype {
	case SYNC:
		fmt.Println("Data base synchronisation, type 1")
	case MAJ, TAKEN, FOLLOW_ON, FOLLOW_OFF:
		fmt.Println(token.Spec.(Ballid).Id)
	case POS:
	case NEW_BALL:
		fmt.Println(token.Spec.(New_ball).Title)
		fmt.Println(token.Spec.(New_ball).Octets)
		fmt.Println(token.Spec.(New_ball).Message)
	case SEND_BALL:
		fmt.Println(token.Spec.(Send_ball).Id)
		fmt.Println(token.Spec.(Send_ball).Octets)
		fmt.Println(token.Spec.(Send_ball).Message)
	case ACK:
		fmt.Println(token.Spec.(Ack).Atype)
		fmt.Println(token.Spec.(Ack).Status)
	}
}
