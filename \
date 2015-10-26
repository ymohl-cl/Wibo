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
	"Wibo/protocol"
	"Wibo/server"
	"container/list"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
func handleConnection(conn net.Conn, Db *sql.DB, Logger *log.Logger, Serv *server.Server) {
	Data := new(answer.Data)
	Data.Lst_req = list.New()
	Data.Lst_asw = list.New()
	Data.Lst_ball = Serv.Lst_ball
	Data.Lst_users = Serv.Lst_users
	Data.Lst_devices = Serv.Lst_Devices
	Data.Lst_work = Serv.Lst_workBall
	Data.Logged = UNKNOWN
	Data.Logger = Logger

	Data.Logger.Println("New connection with (remote Address): ", conn.RemoteAddr())
	defer conn.Close()
	defer Data.Lst_req.Init()
	defer Data.Lst_asw.Init()
	for {
		buff := make([]byte, 1024)
		size, er := conn.Read(buff)
		if er != nil && er == io.EOF {
		} else if er != nil {
			Data.Logger.Printf("Read %d octets, error: %s\n", size, er)
			return
		} else {
			Data.Logger.Printf("Packet received by (remote Address): %s| Detail packet:\n%s\n", conn.RemoteAddr(), buff) // Print Verification

			Token := new(protocol.Request)
			er, er2 := Token.Get_request(buff)
			if er != nil || er2 != nil {
				Data.Logger.Printf("Remote Address: %s| Get_request error1: %s, Get_request_erro2: %s\n", conn.RemoteAddr(), er)
				return
			} else {
				fmt.Println("Received: ") // Print Verification
				Token.Print_token_debug() // Print Verification
			}
			Data.Lst_req.PushBack(Token)
			if Data.Check_lstrequest() == true {
				er = Data.Get_answer(Serv.Tab_wd, Db)
				if er != nil {
					Data.Logger.Printf("Remote Address: %s| Get_answer error: %s\n", conn.RemoteAddr(), er)
					Front := Data.Lst_asw.Front()
					if Front != nil {
						size, er = conn.Write(Front.Value.([]byte))
						Data.Logger.Printf("Remote Address: %s| retour de conn.Write, size: %d, er: %s\n", conn.RemoteAddr(), size, er)
					}
					return
				} else {
					Front := Data.Lst_asw.Front()
					fmt.Println("Answer sending:")    // Print Verification
					fmt.Println(Front.Value.([]byte)) // Print Verification
					size, er = conn.Write(Front.Value.([]byte))
					Data.Logger.Printf("Remote Address: %s| retour de conn.Write, size: %d, er: %s\n", conn.RemoteAddr(), size, er)
					Data.Lst_asw.Remove(Front)
				}
			} else {
				awr := Data.Get_aknowledgement(Data.Lst_users)
				size, er = conn.Write(awr)
				Data.Logger.Printf("Remote Address: %s| retour de conn.Write (exhange multiple packets), size: %d, er: %s\n", conn.RemoteAddr(), size, er)
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
func Listen(Serv *server.Server, Db *sql.DB) {
	var Logger *log.Logger

	file, er := os.Create("Logsfile.txt")
	if er != nil {
		Serv.Logger.Println("os.Create error: ", er)
		os.Exit(-1)
	}
	Logger = log.New(file, "logger: ", log.Llongfile|log.Ldate|log.Ltime)

	ln, er := net.Listen("tcp", ":45899")
	if er != nil {
		Serv.Logger.Println("os.Create error: ", er)
		os.Exit(-1)
	}
	defer ln.Close()

	for {
		conn, er := ln.Accept()
		if er != nil {
			Serv.Logger.Println("Accept error: ", er)
		}
		go handleConnection(conn, Db, Logger, Serv)
	}
}
