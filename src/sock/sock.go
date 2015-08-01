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
	"answer"
	"ballon"
	"container/list"
	"fmt"
	"io"
	"net"
	"owm"
	"protocol"
	"users"
)

/*
** handleConnection received client's requests and manages the exchange with
** client
 */
func handleConnection(conn net.Conn, Lst_users *users.All_users, Lst_ball *ballon.All_ball, Tab_wd *owm.All_data) {
	Data := new(answer.Data)
	Data.Lst_req = list.New()
	Data.Lst_asw = list.New()
	Data.Lst_ball = Lst_ball
	Data.Lst_users = Lst_users

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
			fmt.Println("decode ...")
			Token := new(protocol.Lst_req_sock)
			err := Token.Get_request(buff)
			fmt.Println("Second test")
			fmt.Println(Token)
			if err != nil {
				fmt.Println(err)
				return
			} else {
				/* ! CECI EST POUR FAIRE DES TESTS ! */
				fmt.Println("Receive:")
				Token.Print_token_debug()
				/* FIN DES TESTS */
			}
			Data.Lst_req.PushBack(*Token)
			if Data.Check_lstrequest() == true {
				err = Data.Get_answer(Tab_wd)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					Front := Data.Lst_asw.Front()
					fmt.Println("Answer sending:")
					fmt.Println(Front.Value.([]byte))
					conn.Write(Front.Value.([]byte))
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
func Listen(Lst_users *users.All_users, Lst_ball *ballon.All_ball, Tab_wd *owm.All_data) {
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
		go handleConnection(conn, Lst_users, Lst_ball, Tab_wd)
	}
}
