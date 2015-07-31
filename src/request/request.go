<<<<<<< HEAD
//# ************************************************************************** #
//#                                                                            #
//#                                                       :::      ::::::::    #
//#  request.go                                            :+:      :+:    :+:    #
//#                                                   +:+ +:+         +:+      #
//#  by: ymohl-cl <ymohl-cl@student.42.fr>          +#+  +:+       +#+         #
//#                                               +#+#+#+#+#+   +#+            #
//#  created: 2015/06/11 13:13:33 by ymohl-cl          #+#    #+#              #
//#  updated: 2015/06/11 13:16:35 by ymohl-cl         ###   ########.fr        #
//#                                                                            #
//# ************************************************************************** #

package request

/*
** Ce package permet de recevoir des requetes http.
** ! Ce package n'est encore utilise pour l'echange client serveur !
 */

import (
	"ballon"
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

/*
** Ceci est un test de connection reponse avec le client fonctionnel.
 */
=======
//header

/*
** Package to get all request and give a good answer
 */

package request

import (
	"fmt"
	//	"html"
	"net/http"
	//	"net/url"
	"github.com/Wibo/src/ballon"
	//	"container/list"
	"encoding/json"
)

>>>>>>> testMerge
func Manage_request(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.RequestURI())
	m := req.URL.Query()
	fmt.Println(m)
	fmt.Println(m.Encode())

<<<<<<< HEAD
	// to test with one ball
	tmp_lst := list.New()
	var check_test ballon.Checkpoints
	check_test.Coord.Longitude = 2.316055
	check_test.Coord.Latitude = 48.833086
	check_test.Date = time.Now()

	Lst_test := [2]ballon.Ball{}
	bal1 := ballon.Ball{Name: "toto"}
	bal2 := ballon.Ball{Name: "tata"}

	bal1.Coord = tmp_lst.PushBack(check_test)
	bal2.Coord = bal1.Coord

=======
	Lst_test := [2]ballon.Ball{}
	bal1 := ballon.Ball{Name: "toto", Coord: ballon.Coordinates{Longitude: 10.1, Latitude: 11.1}}
	bal2 := ballon.Ball{Name: "tata", Coord: ballon.Coordinates{Longitude: 42.1, Latitude: 13.1}}
>>>>>>> testMerge
	Lst_test[0] = bal1
	Lst_test[1] = bal2

	test, _ := json.Marshal(Lst_test)
	w.Write([]byte(test))
}

<<<<<<< HEAD
/* Init_handle_request permet d'initialiser le fonctionnement http du serveur */
=======
>>>>>>> testMerge
func Init_handle_request() {
	http.HandleFunc("/", Manage_request)
}
