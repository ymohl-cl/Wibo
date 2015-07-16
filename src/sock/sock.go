// header

package sock

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"users"
)

type Position struct {
	Longitude float32
	Latitude  float32
}

type Id_ballon struct {
	IdBallon int64
}

type Ballon struct {
	Title          string
	Longitude_user float32
	Latitude_user  float32
	Message        string
}

type Lst_req_sock struct {
	NbrOctet int16
	Type     int16
	NbrPack  int32
	NumPack  int32
	IdMobile int64
	Union    interface{}
}

func check_checksum(buff []byte) error {
	return nil
}

func Take_position(TypBuff *bytes.Buffer) (Pos Position, er error) {
	// Get Longitute
	err := binary.Read(TypBuff, binary.BigEndian, &Pos.Longitude)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Pos longitute")
		return Pos, er
	}
	// Get Type Latitude
	err = binary.Read(TypBuff, binary.BigEndian, &Pos.Latitude)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Pos latitude")
		return Pos, er
	}
	return Pos, er
}

func Take_ball(TypBuff *bytes.Buffer) (Id Id_ballon, er error) {
	// Get Type Latitude
	err := binary.Read(TypBuff, binary.BigEndian, &Id.IdBallon)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on ID_ballon")
		return Id, er
	}
	return Id, er
}

func Take_newBall(TypBuff *bytes.Buffer) (Ball Ballon, er error) {
	// Get Title
	var err error
	Ball.Title, err = TypBuff.ReadString(0)
	fmt.Println(len(Ball.Title))
	if len(Ball.Title) == 1 {
		er = errors.New("Add content from socket error, ReadString return error on Ball.Title")
		return Ball, er
	}
	TypBuff.Next((16 - len(Ball.Title)))
	// Get Longitude_user
	err = binary.Read(TypBuff, binary.BigEndian, &Ball.Longitude_user)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Ball longitute")
		return Ball, er
	}
	// Get Latitude_user
	err = binary.Read(TypBuff, binary.BigEndian, &Ball.Latitude_user)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Ball latitude")
		return Ball, er
	}
	TypBuff.Next(8)
	// Get Message
	Ball.Message, err = TypBuff.ReadString(0)
	fmt.Println(len(Ball.Message))
	if 1 == len(Ball.Message) {
		er = errors.New("Add content from socket error, ReadString return error on Ball.Message")
		return Ball, er
	}
	return Ball, er
}

func add_content(buff []byte, user *users.User) (Token Lst_req_sock, er error) {
	TypBuff := bytes.NewBuffer(buff)

	// Get NbrOctet
	err := binary.Read(TypBuff, binary.BigEndian, &Token.NbrOctet)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on NbrOctet")
		return Token, er
	}
	// Get Type request
	err = binary.Read(TypBuff, binary.BigEndian, &Token.Type)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on Type request")
		return Token, er
	}
	TypBuff.Next(4)
	// Get NbrPacket
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NbrPack)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on NbrPack")
		return Token, er
	}
	// Get NumPacket
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NumPack)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on NumPack")
		return Token, er
	}
	// Get IdMobile
	err = binary.Read(TypBuff, binary.BigEndian, &Token.IdMobile)
	if err != nil {
		er = errors.New("Add content from socket error, Binary.Read return error on IdMobile")
		return Token, er
	}
	// Get next content on the buffer
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
	//	if Token.IdMobile != user.Device {
	//		return nil
	//	}
	return Token, nil
}

func Check_finish(Lst_req *list.List) bool {
	Last := Lst_req.Back()

	if Last.Value.(Lst_req_sock).NbrPack == (Last.Value.(Lst_req_sock).NumPack - 1) {
		return true
	} else {
		return false // normalement false ici
	}
}

func Get_answer(Lst_req *list.List) (answer []byte) {
	return answer
}

func Get_aknowledgement(Lst_req *list.List) (answer []byte) {
	return answer
}

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

func handleConnection(conn net.Conn, Lst_users *users.All_users) {
	Lst_req := list.New()
	user := new(users.User)

	defer conn.Close()
	defer Lst_req.Init()
	for {
		buff := make([]byte, 1024)
		_, err := conn.Read(buff) // conn.Read(buff) retourne la taille du buffer et err
		if err != nil {
			return
		}
		answer := buff
		err = check_checksum(buff)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Ajoute le message a la liste.
		//		fmt.Println(buff)
		fmt.Println("New request:")
		Token, err := add_content(buff, user)
		if err != nil {
			fmt.Println(err)
			return
		} else { // debug
			Print_token_debug(Token)
		}
		Lst_req.PushBack(Token)
		// Verifie si la liste est complete
		//si c'est le cas recupere une reponse
		//si c'est pas le cas envoi un acknoldgement.
		if Check_finish(Lst_req) == true {
			//			answer := Get_answer(Lst_req)
			fmt.Println(answer)
			conn.Write(answer)
		} else {
			//			answer := Get_aknowledgement(Lst_req)
			fmt.Println(answer)
			conn.Write(answer)
		}
		answer = nil
		buff = nil
	}
}

func Listen(Lst_users *users.All_users) {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Error listen:", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error Accept:", err)
		}
		go handleConnection(conn, Lst_users)
	}
}
