package request

/*
** Ce package permet de recevoir des requetes http.
** ! Ce package n'est encore utilise pour l'echange client serveur !
 */

import (
	"Wibo/ballon"
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Manage_request(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.RequestURI())
	m := req.URL.Query()
	fmt.Println(m)
	fmt.Println(m.Encode())

	// to test with one ball
	tmp_lst := list.New()
	var check_test ballon.Checkpoint
	check_test.Coord.Lon = 2.316055
	check_test.Coord.Lat = 48.833086
	check_test.Date = time.Now()

	Lst_test := [2]ballon.Ball{}
	bal1 := ballon.Ball{Title: "toto"}
	bal2 := ballon.Ball{Title: "tata"}

	bal1.Coord = tmp_lst.PushBack(check_test)
	bal2.Coord = bal1.Coord

	Lst_test[0] = bal1
	Lst_test[1] = bal2

	test, _ := json.Marshal(Lst_test)
	w.Write([]byte(test))
}

/* Init_handle_request permet d'initialiser le fonctionnement http du serveur */
func Init_handle_request() {
	http.HandleFunc("/", Manage_request)
}
