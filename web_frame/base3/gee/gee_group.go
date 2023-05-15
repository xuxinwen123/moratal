package gee

import (
	"log"
	"net/http"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup //当前分组的parent
	engine      *EngineGroup
}
type EngineGroup struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func NewEngineGroup() *EngineGroup {
	engine := &EngineGroup{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}
func (group *RouterGroup) addRouter(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Router %s-%s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

// Use middleware
// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (e *EngineGroup) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
func (e *EngineGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	ctx := newContext(w, r)
	ctx.handlers = middlewares
	e.router.handlerMiddleware(ctx)
	//分组的实现
	//ctx := newContext(w, r)
	//group.router.handler(ctx)
}
