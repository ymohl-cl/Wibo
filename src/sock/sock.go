// header

package sock

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"net"
	"users"
)

type Lst_req_sock struct {
	NbrOctet int16
	Type     int16
	NbrPack  int
	NumPack  int
	IdMobile int64
	Content  []byte
}

func check_checksum(buff []byte) error {
	return nil
}

func add_content(buff []byte, user *users.User) (Token Lst_req_sock) {
	tmpBuf := buff
	TypBuff := bytes.NewBuffer(buff)

	// Get NbrOctet
	err := binary.Read(TypBuff, binary.BigEndian, &Token.NbrOctet)
	if err != nil {
		//		return nil
	}
	//	Token.NbrOctet -= 24
	tmpBuf = TypBuff.Next(2)
	//	TypBuff = bytes.NewBuffer(buff)
	// Get Type request
	err = binary.Read(TypBuff, binary.BigEndian, &Token.Type)
	if err != nil {
		//		return nil
	}
	tmpBuf = TypBuff.Next(6)
	//	TypBuff = bytes.NewBuffer(buff)
	// Get NbrPacket
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NbrPack)
	if err != nil {
		//		return nil
	}
	tmpBuf = TypBuff.Next(4)
	//	TypBuff = bytes.NewBuffer(buff)
	// Get NumPacket
	err = binary.Read(TypBuff, binary.BigEndian, &Token.NumPack)
	if err != nil {
		//		return nil
	}
	tmpBuf = TypBuff.Next(4)
	//	TypBuff = bytes.NewBuffer(buff)
	// Get IdMobile
	err = binary.Read(TypBuff, binary.BigEndian, &Token.IdMobile)
	if err != nil {
		//		return nil
	}
	tmpBuf = TypBuff.Next(8)
	Token.Content = tmpBuf
	//	if Token.IdMobile != user.Device {
	//		return nil
	//	}
	return Token
}

func Check_finish(Lst_req *list.List) bool {
	Last := Lst_req.Back()
	fmt.Printf("%#v\n", Last)
	fmt.Printf("%#v\n", Last.Value)
	//	fmt.Println(Last.Type.(Lst_req_sock))
	if Last != nil {
		return true
	} else {
		return true // normalement false ici
	}
	//	if Last.NbrPack == Last.NumPack {
	//		return true
	//	} else {
	//		return false
	//	}
	//	return false
}

func Get_answer(Lst_req *list.List) (answer []byte) {
	return answer
}

func Get_aknowledgement(Lst_req *list.List) (answer []byte) {
	return answer
}

func handleConnection(conn net.Conn, Lst_users *users.All_users) {
	Lst_req := list.New()
	user := new(users.User)

	defer conn.Close()
	defer Lst_req.Init()
	for {
		buff := make([]byte, 1024)
		// _ == leng pour les debugs
		leng, err := conn.Read(buff)
		if err != nil {
			return
		}
		err = check_checksum(buff)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Ajoute le message a la liste.
		fmt.Println(buff)
		toto := add_content(buff, user)
		fmt.Println(toto.NbrOctet)
		fmt.Println(toto.Type)
		fmt.Println(toto.NbrPack)
		fmt.Println(toto.NumPack)
		fmt.Println(toto.IdMobile)
		check := Lst_req.PushBack(toto)
		if check == nil {
			fmt.Println("Error decoding data")
			return
		}
		// Verifie si la liste est complete
		//si c'est le cas recupere une reponse
		//si c'est pas le cas envoi un acknoldgement.
		if Check_finish(Lst_req) == true {
			answer := Get_answer(Lst_req)
			conn.Write(answer)
		} else {
			answer := Get_aknowledgement(Lst_req)
			conn.Write(answer)
		}
		str := string(buff[:leng])
		fmt.Println("Message Received:", str)
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
