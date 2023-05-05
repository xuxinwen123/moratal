package common

import (
	"github.com/sirupsen/logrus"
	"moratal/apirouter/util"
	"strings"
)

// NewRouter 创建路由器
func NewRouter(api string) (router *Router) {
	router = newRouter()
	if api != "" {
		router.AddRouter(parseApi(api)...)
	}
	return
}

type Router struct {
	roots   map[string]*segment
	routers map[string]map[string]bool // 路由缓存, key: method, value: map[pattern]bool, 在添加路由前先查询缓存中是否存在该路由
}
type ApiInfo struct {
	Method  string
	Pattern string
}

func newRouter() *Router {
	return &Router{
		roots:   make(map[string]*segment),
		routers: make(map[string]map[string]bool),
	}
}
func (r *Router) AddRouter(apiList ...ApiInfo) {
	for _, api := range apiList {
		api.Method = strings.ToUpper(api.Method)
		// check if route is already exist, if not in then add it to cache
		if r.inCache(api.Method, api.Pattern) {
			logrus.Errorf("req url :{method:%s,pattern:%s} is alreay registered", api.Method, api.Pattern)
			continue
		}
		//add new router
		parts := ParsePattern(api.Pattern)
		if _, ok := r.roots[api.Method]; !ok {
			r.roots[api.Method] = &segment{}
		}
		r.roots[api.Method].insert(api.Method, api.Pattern, parts, 0)
		logrus.Debugf("register api:'[%s] %s'", api.Method, api.Pattern)

		r.addCache(api.Method, api.Pattern)
	}
	//排序：静态路由>动态路由
	r.sortRouters()
}

/*
ParsePattern

	e.g. pattern='/api/nodes/{id:int}' -> ["api", "nodes", "{id:int}"]
	e.g. pattern='/api/nodes/3' -> ["api", "nodes", "3"]
*/
func ParsePattern(pattern string) []string {
	if pattern == "" {
		return nil
	}
	// pattern = "/api/user/new" --> (remove the first '/') "api/user/new" --> (split) ["api", "user", "new"]
	// pattern = "/api/user/{id:int}" --> (remove the first '/') "api/user/{id:int}" --> (split) ["api", "user", "{id:int}"]
	return strings.Split(pattern[1:], "/")
}
func (r *Router) inCache(method, pattern string) bool {
	if pattern == "" {
		return false
	}
	m, ok := r.routers[method]
	if !ok {
		return false
	}
	return m[pattern]
}
func (r *Router) addCache(method, patten string) {
	if patten == "" {
		return
	}
	rs, ok := r.routers[method]
	if !ok {
		r.routers[method] = make(map[string]bool)
		r.routers[method][patten] = true
		return
	}
	rs[patten] = true
}
func (r *Router) sortRouters() {
	for _, rootSegment := range r.roots {
		rootSegment.sortChildren()
	}
}
func parseApi(apiConf string) (apiList []ApiInfo) {
	apiMapping := make(map[string]string)
	if err := util.UnmarshalFromFile(&apiMapping, apiConf, ""); err != nil {
		panic(err)
	}
	// eg:/api/users/getUserInfo    GET,POST
	//为了解决method不同，但是pattern相同的问题，所以split methods
	for pattern, methods := range apiMapping {
		splitMethod := strings.Split(methods, ",")
		for i := 0; i < len(splitMethod); i++ {
			apiList = append(apiList, ApiInfo{Method: splitMethod[0], Pattern: pattern})
		}
	}
	return
}
