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
	"net"
	"protocol"
	"users"
)

/*
** handleConnection recoit une requete, lance le traitement de cette requete
** Ensuite il ecrit la reponse en retour du traitement pour l'envoyer
** et ecoute a nouveau le client.
** conn.Read(buff) retourne la taille du buff et error
 */
func handleConnection(conn net.Conn, Lst_users *users.All_users, Lst_ball *ballon.All_ball) {
	Data := new(answer.Data)
	Data.Lst_req = list.New()
	Data.Lst_asw = list.New()
	Data.Lst_ball = Lst_ball
	Data.Lst_users = Lst_users

	defer conn.Close()
	defer Data.Lst_req.Init()
	defer Data.Lst_asw.Init()
	for {
		buff := make([]byte, 1024)
		_, err := conn.Read(buff)
		if err != nil {
			return
		}
		fmt.Println("New request:")
		Token, err := protocol.Add_content(buff)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			/* ! CECI EST POUR FAIRE DES TESTS ! */
			protocol.Print_token_debug(Token)
			/* FIN DES TESTS */
		}
		fmt.Println(Token)
		Data.Lst_req.PushBack(Token)
		if answer.Check_packets_list(Data.Lst_req.Front()) == true {
			fmt.Println("check finish: ok")
			err = Data.Get_answer()
			if err != nil {
				fmt.Println("Erreur Data.Get_answer")
				fmt.Println(err)
			} else {
				fmt.Println("Packet found and send")
				Front := Data.Lst_asw.Front()
				if Front != nil {
					fmt.Println("Front == nil")
					fmt.Println("exit")
				} else {
					fmt.Println("Front != nil")
					fmt.Println(Front)
				}
				conn.Write(Front.Value.([]byte))
				Data.Lst_asw.Remove(Front)
			}
		} else {
			fmt.Println("Check finish: ko")
			awr, err := answer.Get_aknowledgement(Data.Lst_req, Lst_users)
			if err != nil {
				fmt.Println(err)
			} else {
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
func Listen(Lst_users *users.All_users, Lst_ball *ballon.All_ball) {
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
		go handleConnection(conn, Lst_users, Lst_ball)
	}
}
