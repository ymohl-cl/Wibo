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
	Lonuser float64
	Latuser float64
	Message string
}

//iddevice int64 // Deviendra une string ou un buffer ..
type Request struct {
	Octets   int16
	Rtype    int16
	Nbrpck   int32
	Numpck   int32
	Deviceid int64
	Spec     interface{}
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

/* Decode type position */
func Request_position(TypBuff *bytes.Buffer) (pos Position, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &pos.Lon)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return pos, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &pos.Lat)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return pos, er
	}
	fmt.Println("Value dans take Position")
	fmt.Println(pos)
	return pos, er
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
	if len(ball.Title) == 1 || err != nil {
		er = errors.New("Get_newball in protocol: Error ReadString")
		return ball, er
	}
	TypBuff.Next((16 - len(ball.Title)))
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Lonuser)
	if err != nil {
		er = errors.New("Get_newball in protocol: Error binary.Read")
		return ball, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Latuser)
	if err != nil {
		er = errors.New("Get_newball in protocol: Error binary.Read")
		return ball, er
	}
	TypBuff.Next(8)
	ball.Message, err = TypBuff.ReadString(0)
	if 1 == len(ball.Message) {
		er = errors.New("Get_newball in protocol: Error ReadString")
		return ball, er
	}
	return ball, er
}

/* Decode type sendball  */
func Request_sendball(TypBuff *bytes.Buffer) (ball New_ball, er error) {
	return ball, er
}

/*
** Master function that tokenizes client's request.
** See the protocol
 */
func (token *Request) Get_request(buff []byte) (er error) {
	TypBuff := bytes.NewBuffer(buff)

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
	err = binary.Read(TypBuff, binary.BigEndian, &token.Deviceid)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	TypBuff.Next(32)
	switch token.Rtype {
	case SYNC:
		return er
	case MAJ, TAKEN, FOLLOW_ON, FOLLOW_OFF:
		token.Spec, er = Request_idball(TypBuff)
	case POS:
		token.Spec, er = Request_position(TypBuff)
	case NEW_BALL:
		token.Spec, er = Request_newball(TypBuff)
	case SEND_BALL:
		token.Spec, er = Request_sendball(TypBuff)
	case ACK:
		token.Spec, er = Request_ack(TypBuff)
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
	fmt.Println(token.Deviceid)
	fmt.Println("Type request:")
	switch token.Rtype {
	case SYNC:
		fmt.Println("Data base synchronisation, type 1")
	case MAJ, TAKEN, FOLLOW_ON, FOLLOW_OFF:
		fmt.Println(token.Spec.(Ballid).Id)
	case POS:
		fmt.Println(token.Spec.(Position).Lon)
		fmt.Println(token.Spec.(Position).Lat)
	case NEW_BALL:
		fmt.Println(token.Spec.(New_ball).Title)
		fmt.Println(token.Spec.(New_ball).Lonuser)
		fmt.Println(token.Spec.(New_ball).Latuser)
		fmt.Println(token.Spec.(New_ball).Message)
	case SEND_BALL:
		fmt.Println("Renvoi un ballon non gere pour le moment")
	case ACK:
		fmt.Println(token.Spec.(Ack).Atype)
		fmt.Println(token.Spec.(Ack).Status)
	}
}
