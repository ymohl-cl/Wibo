//header

/*
** Package to get all request and give a good answer
 */

package request

import (
	"fmt"
	//	"html"
	"net/http"
	"net/url"
)

func Manage_request(w http.ResponseWriter, req *http.Request) {
	/*	s := html.EscapeString(req.URL.Path)
		fmt.Println(req.URL.Path)*/
	ul, _ := url.Parse(req.URL.Path)
	fmt.Println(ul.RawQuery)
	fmt.Println(url.ParseQuery(ul.RawQuery))
	w.Write([]byte("<p>BONJOUR !!!! :D</p>"))
}

func Init_handle_request() {
	http.HandleFunc("/", Manage_request)
}
