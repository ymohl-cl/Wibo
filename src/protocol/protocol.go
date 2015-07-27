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

/* Take_position est le decodeur de type 1 */
func Take_position(TypBuff *bytes.Buffer) (Pos Position, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &Pos.Longitude)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Pos longitute")
		return Pos, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Pos.Latitude)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Pos latitude")
		return Pos, er
	}
	return Pos, er
}

/* Take_ball est le decodeur de type 2, 3, et 4 */
func Take_ball(TypBuff *bytes.Buffer) (Id Id_ballon, er error) {
	err := binary.Read(TypBuff, binary.BigEndian, &Id.IdBallon)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on ID_ballon")
		return Id, er
	}
	return Id, er
}

/* TaKe_newBall est le decodeur de type 5 */
func Take_newBall(TypBuff *bytes.Buffer) (Ball Ballon, er error) {
	var err error
	Ball.Title, err = TypBuff.ReadString(0)
	fmt.Println(len(Ball.Title))
	if len(Ball.Title) == 1 {
		er = errors.New("Add content from socket error, ReadString return error on Ball.Title")
		return Ball, er
	}
	TypBuff.Next((16 - len(Ball.Title)))
	err = binary.Read(TypBuff, binary.BigEndian, &Ball.Longitude_user)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Ball longitute")
		return Ball, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Ball.Latitude_user)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Ball latitude")
		return Ball, er
	}
	TypBuff.Next(8)
	Ball.Message, err = TypBuff.ReadString(0)
	if 1 == len(Ball.Message) {
		er = errors.New("Add content from socket error, ReadString return error on Ball.Message")
		return Ball, er
	}
	return Ball, er
}

/*
** Add_content recupere le contenu du header de la requete recu dans buff et
** l'analyse pour creer une requete exploitable par le serveur, en appelant
** le decodeur du type specifier dans le header (Voir protocole Wibo sur trello)
 */
func Add_content(buff []byte) (Token Lst_req_sock, er error) {
	TypBuff := bytes.NewBuffer(buff)

	err := binary.Read(TypBuff, binary.BigEndian, &Token.NbrOctet)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on NbrOctet")
		return Token, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Token.Type)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Type request")
		return Token, er
	}
	TypBuff.Next(4)
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NbrPack)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on NbrPack")
		return Token, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NumPack)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on NumPack")
		return Token, er
	}
	err = binary.Read(TypBuff, binary.BigEndian, &Token.IdMobile)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on IdMobile")
		return Token, er
	}
	TypBuff.Next(32)
	switch Token.Type {
	case 1:
		Token.Union, er = Take_position(TypBuff)
		if er != nil {
			return Token, er
		}
	case 2, 3, 4:
		Token.Union, er = Take_ball(TypBuff)
		if er != nil {
			return Token, er
		}
	case 5:
		Token.Union, er = Take_newBall(TypBuff)
		if er != nil {
			return Token, er
		}
	}
	return Token, nil
}

/* Fonction pour faire des prints de debug sur une requete recu */
func Print_token_debug(Token Lst_req_sock) {
	fmt.Println(Token.NbrOctet)
	fmt.Println(Token.Type)
	fmt.Println(Token.NbrPack)
	fmt.Println(Token.NumPack)
	fmt.Println(Token.IdMobile)
	switch Token.Type {
	case 1:
		fmt.Println(Token.Union.(Position).Longitude)
		fmt.Println(Token.Union.(Position).Latitude)
	case 2, 3, 4:
		fmt.Println(Token.Union.(Id_ballon).IdBallon)
	case 5:
		fmt.Println(Token.Union.(Ballon).Title)
		fmt.Println(Token.Union.(Ballon).Longitude_user)
		fmt.Println(Token.Union.(Ballon).Latitude_user)
		fmt.Println(Token.Union.(Ballon).Message)
	}
}
