package gosaas

import (
	"path"
	"strings"
)

// ShiftPath splits the request URL head and tail.
//
// This is useful to perform routing inside the ServeHTTP function.ShiftPath
//
// Example usage:
//
// 	package yourapp
//
// 	import (
// 		"net/http"
// 		"github.com/dstpierre/gosaas"
// 	)
//
// 	func main() {
// 		routes := make(map[string]*gosaas.Route)
// 		routes["speak"] = &gosaas.Route{Handler: speak{}}
// 		mux := gosaas.NewServer(routes)
// 		http.ListenAndServe(":8080", mux)
// 	}
//
// 	type speak struct{}
// 	func (s speak) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 		var head string
// 		head, r.URL.Path = gosaas.ShiftPath(r.URL.Path)
// 		if head == "loud" {
// 			s.scream(w, r)
// 		} else {
// 			s.whisper(w, r)
// 		}
// 	}
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
