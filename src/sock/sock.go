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
	//	"answer"
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
func handleConnection(conn net.Conn, Lst_users *users.All_users) {
	Lst_req := list.New()
	user := new(users.User)

	defer conn.Close()
	defer Lst_req.Init()
	for {
		buff := make([]byte, 1024)
		_, err := conn.Read(buff)
		if err != nil {
			return
		}
		fmt.Println("New request:")
		answer := buff
		Token, err := protocol.Add_content(buff, user)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			/* ! CECI EST POUR FAIRE DES TESTS ! */
			protocol.Print_token_debug(Token)
			/* FIN DES TESTS */
		}
		Lst_req.PushBack(Token)
		//		if protocol.Check_finish(Lst_req) == true {
		//			awr, err := answer.Get_answer(Lst_req, Lst_users)
		//			if err != nil {
		//				fmt.Println(err)
		//			} else {
		//				conn.Write(awr)
		//			}
		//		} else {
		//			awr, err := answer.Get_aknowledgement(Lst_req, Lst_users)
		//			if err != nil {
		//				fmt.Println(err)
		//			} else {
		//				conn.Write(awr)
		conn.Write(answer)
		//			}
		//		}
		buff = nil
	}
}

/*
** Listen va ecouter les connections entrante sur le port 8081
** Elle va accepter une demande de connection et lancer le handleConnection
** handleConnection va recuperer et repondre au requete du client jusqu'a
** arriver a un etat close.
 */
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
