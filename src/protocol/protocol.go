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
	Type int16
	Flag int16
}

type Position struct {
	Longitude float64
	Latitude  float64
}

type Id_ballon struct {
	IdBallon int64
}

type Ballon struct {
	Title          string
	Longitude_user float64
	Latitude_user  float64
	Message        string
}

type Lst_req_sock struct {
	NbrOctet int16
	Type     int16
	NbrPack  int32
	NumPack  int32
	IdMobile int64 // Deviendra une string ou un buffer ..
	Union    interface{}
}

/* Decode type Ack */
func Request_ack(TypBuff *bytes.Buffer) (ack Ack, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &ack.Type)
	if err != nil {
		er = errors.New("Get_ack in protocol: Error binary.Read")
		return ack, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &ack.Flag)
	if err != nil {
		er = errors.New("Get_ack in protocol: Error binary.Read")
		return ack, er
	}
	return ack, er
}

/* Decode type position */
func Request_position(TypBuff *bytes.Buffer) (pos Position, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &pos.Longitude)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return pos, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &pos.Latitude)
	if err != nil {
		er = errors.New("Get_position in protocol: Error binary.Read")
		return pos, er
	}
	fmt.Println("Value dans take Position")
	fmt.Println(pos)
	return pos, er
}

/* Decode types MAJ, TAKEN, FOLLOW_ON, and FOLLOW_OFF */
func Request_idball(TypBuff *bytes.Buffer) (id Id_ballon, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &id.IdBallon)

	if err != nil {
		er = errors.New("Get_idball in protocol: Error binary.Read")
		return id, er
	}
	return id, er
}

/* Decode type newball */
func Request_newball(TypBuff *bytes.Buffer) (ball Ballon, er error) {
	var err error

	ball.Title, err = TypBuff.ReadString(0)
	if len(ball.Title) == 1 || err != nil {
		er = errors.New("Get_newball in protocol: Error ReadString")
		return ball, er
	}
	TypBuff.Next((16 - len(ball.Title)))
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Longitude_user)
	if err != nil {
		er = errors.New("Get_newball in protocol: Error binary.Read")
		return ball, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &ball.Latitude_user)
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
func Request_sendball(TypBuff *bytes.Buffer) (ball Ballon, er error) {
	return ball, er
}

/*
** Master function that tokenizes client's request.
** See the protocol
 */
func (Token *Lst_req_sock) Get_request(buff []byte) (er error) {
	TypBuff := bytes.NewBuffer(buff)

	er = nil
	err := binary.Read(TypBuff, binary.BigEndian, &Token.NbrOctet)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Token.Type)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	TypBuff.Next(4)
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NbrPack)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NumPack)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Token.IdMobile)
	if err != nil {
		er = errors.New("Get_request in protocol: Error binary.Read")
		return er
	}
	TypBuff.Next(32)
	switch Token.Type {
	case SYNC:
		return er
	case MAJ, TAKEN, FOLLOW_ON, FOLLOW_OFF:
		Token.Union, er = Request_idball(TypBuff)
	case POS:
		Token.Union, er = Request_position(TypBuff)
	case NEW_BALL:
		Token.Union, er = Request_newball(TypBuff)
	case SEND_BALL:
		Token.Union, er = Request_sendball(TypBuff)
	case ACK:
		Token.Union, er = Request_ack(TypBuff)
	}
	if er != nil {
		return er
	}
	fmt.Println("First test")
	fmt.Println(Token)
	return er
}

/* Print request to debug */
func (Token *Lst_req_sock) Print_token_debug() {
	fmt.Println("Request header:")
	fmt.Println(Token.NbrOctet)
	fmt.Println(Token.Type)
	fmt.Println(Token.NbrPack)
	fmt.Println(Token.NumPack)
	fmt.Println(Token.IdMobile)
	fmt.Println("Type request:")
	switch Token.Type {
	case SYNC:
		fmt.Println("Data base synchronisation, type 1")
	case MAJ, TAKEN, FOLLOW_ON, FOLLOW_OFF:
		fmt.Println(Token.Union.(Id_ballon).IdBallon)
	case POS:
		fmt.Println(Token.Union.(Position).Longitude)
		fmt.Println(Token.Union.(Position).Latitude)
	case NEW_BALL:
		fmt.Println(Token.Union.(Ballon).Title)
		fmt.Println(Token.Union.(Ballon).Longitude_user)
		fmt.Println(Token.Union.(Ballon).Latitude_user)
		fmt.Println(Token.Union.(Ballon).Message)
	case SEND_BALL:
		fmt.Println("Renvoi un ballon non gere pour le moment")
	case ACK:
		fmt.Println(Token.Union.(Ack).Type)
		fmt.Println(Token.Union.(Ack).Flag)
	}
}
