package gee

//
//import (
//	"fmt"
//	"net/http"
//)
//
//type HandlerFunc func(w http.ResponseWriter, r *http.Request)
//type Engine struct {
//	router map[string]HandlerFunc
//}
//
//// New is constructor of gee.Engine
//func New() *Engine {
//	return &Engine{router: make(map[string]HandlerFunc)}
//}
//func (e *Engine) addRouter(method, pattern string, handler HandlerFunc) {
//	key := method + "-" + pattern
//	e.router[key] = handler
//}
//func (e *Engine) GET(pattern string, handler HandlerFunc) {
//	e.addRouter("GET", pattern, handler)
//}
//func (e *Engine) POST(pattern string, handler HandlerFunc) {
//	e.addRouter("POST", pattern, handler)
//}
//func (e *Engine) Run(addr string) (err error) {
//	return http.ListenAndServe(addr, e)
//}
//func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	key := req.Method + "-" + req.URL.Path
//	if handler, ok := e.router[key]; ok {
//		handler(w, req)
//	} else {
//		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
//	}
//}
