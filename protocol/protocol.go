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
	"strings"
)

const (
	_          = iota
	ACK        = 32767
	SYNC       = 1
	MAJ        = 2
	POS        = 3
	TAKEN      = 4
	FOLLOW_ON  = 5
	FOLLOW_OFF = 6
	NEW_BALL   = 7
	SEND_BALL  = 8
	MAGNET     = 9
	// Itinerary is depreacated
	WORKBALL      = 10
	TYPELOG       = 11
	CREATEACCOUNT = 12
	SYNCROACCOUNT = 13
	DELOG         = 14
	STATSUSER     = 15
	STATSBALL     = 16
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

type Taken struct {
	Id         int64
	FlagMagnet int16
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
	Pswd  []byte
}

type Request struct {
	Octets   int16
	Rtype    int16
	Nbrpck   int32
	Numpck   int32
	IdMobile string
	Coord    Position
	Spec     interface{}
}

/* Decode type Ack */
func (tkn *Request) Request_ack(TypBuff *bytes.Buffer) (err error, er error) {
	var ack Ack

	err = binary.Read(TypBuff, binary.BigEndian, &ack.Atype)
	if err != nil {
		er = errors.New("Read ack.Atype")
		return
	}
	err = binary.Read(TypBuff, binary.BigEndian, &ack.Status)
	if err != nil {
		er = errors.New("Read acl.Status")
		return
	}
	tkn.Spec = ack
	return
}

func (tkn *Request) Request_taken(TypBuff *bytes.Buffer) (err error, er error) {
	var take Taken

	err = binary.Read(TypBuff, binary.BigEndian, &take.Id)
	if err != nil {
		er = errors.New("Read take id")
		return
	}
	err = binary.Read(TypBuff, binary.BigEndian, &take.FlagMagnet)
	if err != nil {
		er = errors.New("Read FlagMagnet")
		return
	}
	TypBuff.Next(6)
	tkn.Spec = take
	return
}

/* Decode types MAJ, TAKEN, FOLLOW_ON, and FOLLOW_OFF */
func (tkn *Request) Request_idball(TypBuff *bytes.Buffer) (err error, er error) {
	var id Ballid

	err = binary.Read(TypBuff, binary.BigEndian, &id.Id)
	if err != nil {
		er = errors.New("Read id.Id")
		return
	}
	tkn.Spec = id
	return
}

/* Decode type newball */
func (tkn *Request) Request_newball(TypBuff *bytes.Buffer) (err error, er error) {
	var ball New_ball

	ball.Title, err = TypBuff.ReadString(0)
	if len(ball.Title) == 1 {
		er = errors.New("Read ball.Title")
		return
	}
	TypBuff.Next((16 - len(ball.Title)))
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Octets)
	if err != nil {
		er = errors.New("Read ball.Octets")
		return
	}
	TypBuff.Next(4)
	ball.Message, err = TypBuff.ReadString(0)
	if 1 == len(ball.Message) {
		er = errors.New("Read ball.Message")
		return
	}
	tkn.Spec = ball
	return
}

/* Decode type sendball  */
func (tkn *Request) Request_sendball(TypBuff *bytes.Buffer) (err error, er error) {
	var ball Send_ball

	err = binary.Read(TypBuff, binary.BigEndian, &ball.Id)
	if err != nil {
		er = errors.New("Read ball Id")
		return
	}
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Octets)
	if err != nil {
		er = errors.New("Read ball.Octets")
		return
	}
	TypBuff.Next(4)
	ball.Message, err = TypBuff.ReadString(0)
	if 1 == len(ball.Message) {
		er = errors.New("Read ball Message")
		return
	}
	tkn.Spec = ball
	return
}

func (tkn *Request) Request_Log(TypBuff *bytes.Buffer) (err error, er error) {
	var email [320]byte
	var pswd [64]byte
	var log Log

	err = binary.Read(TypBuff, binary.BigEndian, &email)
	if err != nil {
		er = errors.New("Read email")
		return
	}
	log.Email = string(email[:320])
	log.Email = strings.Trim(log.Email, "\x00")
	err = binary.Read(TypBuff, binary.BigEndian, &pswd)
	if err != nil {
		er = errors.New("Read &pswd")
		return
	}
	log.Pswd = pswd[:64]
	tkn.Spec = log
	return
}

/*
** Master function that tokenizes client's request.
** See the protocol
 */
func (token *Request) Get_request(buff []byte) (err error, er error) {
	TypBuff := bytes.NewBuffer(buff)
	var IdMobile [40]byte
	er = nil

	err = binary.Read(TypBuff, binary.BigEndian, &token.Octets)
	if err != nil {
		er = errors.New("Read token.Octets")
		return err, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &token.Rtype)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return err, er
	}
	TypBuff.Next(4)
	err = binary.Read(TypBuff, binary.BigEndian, &token.Nbrpck)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return err, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &token.Numpck)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return err, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &IdMobile)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return err, er
	}
	token.IdMobile = string(IdMobile[:40])
	token.IdMobile = strings.Trim(token.IdMobile, "\x00")
	err = binary.Read(TypBuff, binary.BigEndian, &token.Coord.Lon)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return err, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &token.Coord.Lat)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return err, er
	}
	switch token.Rtype {
	default:
		er = errors.New("Request type is unknown")
		return err, er
	case SYNC:
		return err, er
	case MAJ, FOLLOW_ON, FOLLOW_OFF, STATSBALL:
		err, er = token.Request_idball(TypBuff)
	case POS, WORKBALL, MAGNET:
	case TAKEN:
		err, er = token.Request_taken(TypBuff)
	case NEW_BALL:
		err, er = token.Request_newball(TypBuff)
	case SEND_BALL:
		err, er = token.Request_sendball(TypBuff)
	case ACK:
		err, er = token.Request_ack(TypBuff)
	case TYPELOG, CREATEACCOUNT:
		err, er = token.Request_Log(TypBuff)
	case SYNCROACCOUNT:
	case DELOG:
	case STATSUSER:
	}
	return err, er
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
	case MAJ, FOLLOW_ON, FOLLOW_OFF, STATSBALL:
		fmt.Println(token.Spec.(Ballid).Id)
	case POS, WORKBALL, MAGNET:
	case TAKEN:
		fmt.Println(token.Spec.(Taken).Id)
		fmt.Println(token.Spec.(Taken).FlagMagnet)
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
	case STATSUSER:
	}
}
