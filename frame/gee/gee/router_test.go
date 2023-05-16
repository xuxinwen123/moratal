package gee

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "/hello/b/c", nil)
	r.addRouter("GET", "/hi/:name", nil)
	r.addRouter("GET", "/assets/*filepath", nil)
	return r
}
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/name"), []string{"p", "name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		log.Fatal("test parse failed")
	}
}
func TestGetRouter(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRouter("GET", "/hello/geektutu")
	if n == nil {
		t.Fatal("get router failed")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("not match /hello/:name")
	}
	if ps["name"] != "geektutu" {
		t.Fatal("name need be equal geektutu")
	}
	fmt.Printf("matched path :%s,params['name']:%s\n", n.pattern, ps["name"])
}
