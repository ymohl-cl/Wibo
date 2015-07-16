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

func Manage_request(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.RequestURI())
	m := req.URL.Query()
	fmt.Println(m)
	fmt.Println(m.Encode())

	Lst_test := [2]ballon.Ball{}
	bal1 := ballon.Ball{Name: "toto", Coord: ballon.Coordinates{Longitude: 10.1, Latitude: 11.1}}
	bal2 := ballon.Ball{Name: "tata", Coord: ballon.Coordinates{Longitude: 42.1, Latitude: 13.1}}
	Lst_test[0] = bal1
	Lst_test[1] = bal2

	test, _ := json.Marshal(Lst_test)
	w.Write([]byte(test))
}

func Init_handle_request() {
	http.HandleFunc("/", Manage_request)
}
