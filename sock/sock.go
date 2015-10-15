//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  sock.go                                            :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  by: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

package sock

import (
	"Wibo/answer"
	"Wibo/ballon"
	"Wibo/owm"
	"Wibo/protocol"
	"Wibo/users"
	"container/list"
	"database/sql"
	"fmt"
	"io"
	"net"
)

const (
	TYPELOG     = 0
	UNKNOWN     = 1
	DEFAULTUSER = 2
	USERLOGGED  = 3
)

/*
** handleConnection received client's requests and manages the exchange with
** client
 */
func handleConnection(conn net.Conn, Lst_users *users.All_users, Lst_ball *ballon.All_ball, Tab_wd *owm.All_data, Db *sql.DB) {
	Data := new(answer.Data)
	Data.Lst_req = list.New()
	Data.Lst_asw = list.New()
	Data.Lst_ball = Lst_ball
	Data.Lst_users = Lst_users
	Data.Logged = UNKNOWN

	fmt.Println("Start handle Connection")
	defer conn.Close()
	defer Data.Lst_req.Init()
	defer Data.Lst_asw.Init()
	for {
		buff := make([]byte, 1024)
		size, err := conn.Read(buff)
		if err != nil && err == io.EOF {
		} else if err != nil {
			fmt.Printf("Read Error, size: %d", size)
			return
		} else {
			fmt.Println("Received")
			fmt.Println(buff)
			Token := new(protocol.Request)
			err := Token.Get_request(buff)
			if err != nil {
				fmt.Println(err)
				return
			} else {
				/* ! CECI EST POUR FAIRE DES TESTS ! */
				fmt.Println("Receive:")
				Token.Print_token_debug()
				/* FIN DES TESTS */
			}
			Etoken := Data.Lst_req.PushBack(Token)
			// Cette partie doit etre dans Answer.
			//			if Data.Logged == UNKNOWN {
			//				if Token.Rtype != TYPELOG {
			//					fmt.Println("Reception anormale")
			//					// return packet negatif.
			//					return
			//				}
			//				Data.Device, err = Dlist.GetDevice(Etoken, Db)
			//				if err != nil {
			//					fmt.Println("Error on GetDevice")
			//					// Return packet negatif.
			//					return
			//				}
			//				Data.User, err = Lst_users.Check_user(Etoken, Db)
			//				if err != nil {
			//					fmt.Println("Error on check users")
			//					// Return packet negatif.
			//					return
			//				}
			//				//Send packet acknowledgement positif
			//			}
			if Data.Check_lstrequest() == true {
				err = Data.Get_answer(Tab_wd, Db)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					Front := Data.Lst_asw.Front()
					fmt.Println("Answer sending:")
					fmt.Println(Front.Value.([]byte))
					size, err = conn.Write(Front.Value.([]byte))
					Data.Lst_asw.Remove(Front)
				}
			} else {
				fmt.Println("Multiple packets exchange is not finish")
				awr := Data.Get_aknowledgement(Lst_users)
				fmt.Println("Answer sending:")
				fmt.Println(awr)
				conn.Write(awr)
			}
		}
		buff = nil
	}
}

/*
** Listen va ecouter les connections entrante sur le port 8081
** Elle va accepter une demande de connection et lancer le handleConnection
** handleConnection va recuperer et repondre au requete du client jusqu'a
** arriver a un etat close.
 */
func Listen(Lst_users *users.All_users, Lst_ball *ballon.All_ball, Tab_wd *owm.All_data, Db *sql.DB) {
	ln, err := net.Listen("tcp", ":45899")
	if err != nil {
		fmt.Println("Error listen:", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error Accept:", err)
		}
		go handleConnection(conn, Lst_users, Lst_ball, Tab_wd, Db)
	}
}
