// header

package sock

import (
	"answer"
	"container/list"
	"fmt"
	"net"
	"protocol"
	"users"
)

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
		//answer := buff
		// Ajoute le message a la liste.
		//		fmt.Println(buff)
		fmt.Println("New request:")
		Token, err := protocol.Add_content(buff, user)
		if err != nil {
			fmt.Println(err)
			return
		} else { // debug
			protocol.Print_token_debug(Token)
		}
		Lst_req.PushBack(Token)
		// Verifie si la liste est complete
		//si c'est le cas recupere une reponse
		//si c'est pas le cas envoi un acknoldgement.
		if protocol.Check_finish(Lst_req) == true {
			awr, err := answer.Get_answer(Lst_req, Lst_users)
			if err != nil {
				fmt.Println(err)
			} else {
				conn.Write(awr)
				//			fmt.Println(answer)
			}
		} else {
			awr, err := answer.Get_aknowledgement(Lst_req, Lst_users)
			if err != nil {
				fmt.Println(err)
			} else {
				conn.Write(awr)
				//			fmt.Println(answer)
			}
		}
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
