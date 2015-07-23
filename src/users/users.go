//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  users.go                                           :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  by: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

package users

import (
	"container/list"
	"time"
)

/*
** Date est la date a laquelle la requete a ete effectue.
** Type_req_client et le type de requete effectue.
 */
type History struct {
	Date            time.Time
	Type_req_client int16
}

/*
** IdMobile est l'identifiant unique du mobile.
** Pour le moment le format exact de l'IdMobile est inconnu.
** History_req est une liste qui sera l'historique des requetes du client
** depuis ce device.
 */
type Device struct {
	IdMobile    int64
	History_req *list.List
}

/*
** Device est la liste des devices de l'utilisateur
** Log est la date de dernier signe de vie utilisateur
 */
type User struct {
	Device *list.List
	Log    time.Time
}

type All_users struct {
	Lst_users *list.List
}

/* Definis si l'utilisateur est considere en ligne ou pas avec un timeout de 2 min */
func (User *User) User_is_online() bool {
	t_now := time.Now()
	t_user := User.Log
	if t_user.Hour() == t_now.Hour() && t_user.Minute() > t_now.Minute()-2 {
		return true
	} else {
		return false
	}
}

/*
** Del_user va supprimer un user directement dans la base de donnee et dans sa propre liste,
** un utilisateur passe en parametre.
 */
func (Lst_users *All_users) Del_user(del_user *User) {
	return
}

/*
** Add_new_user va rajouter directement dans la base de donnee, un utilisateur
** passer en parametre.
 */
func (Lst_users *All_users) Add_new_user(new_user *User) {
	return
}

/* Print_users print tous les utilisateurs */
func (Lst_users *All_users) Print_users() {
	return
}

/*
** Get_users va recuperer tous les utilisateurs dans la base de donnee.
 */
func (Lst_users *All_users) Get_users() error {
	return nil
}
