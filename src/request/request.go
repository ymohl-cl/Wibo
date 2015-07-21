//header

/*
** Package to get all request and give a good answer
 */

package request

import (
	"ballon"
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
)

func Manage_request(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.RequestURI())
	m := req.URL.Query()
	fmt.Println(m)
	fmt.Println(m.Encode())

	// to test with one ball
	tmp_lst := list.New()
	var check_test ballon.Checkpoints
	check_test.Coord.Longitude = 2.316055
	check_test.Coord.Latitude = 48.833086
	check_test.Date = 0

	Lst_test := [2]ballon.Ball{}
	bal1 := ballon.Ball{Name: "toto"}
	bal2 := ballon.Ball{Name: "tata"}

	bal1.Coord = tmp_lst.PushBack(check_test)
	bal2.Coord = bal1.Coord

	Lst_test[0] = bal1
	Lst_test[1] = bal2

	test, _ := json.Marshal(Lst_test)
	w.Write([]byte(test))
}

func Init_handle_request() {
	http.HandleFunc("/", Manage_request)
}
