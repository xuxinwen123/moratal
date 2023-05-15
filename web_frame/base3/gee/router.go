package gee

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg: roots["GET"] roots["POST]
// handles key eg,handlers['GET-/p/:lang/doc'] handlers['POST-/p/book']
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	fmt.Println("the parts is", parts)
	return parts
}

func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	//找到匹配的路由
	search := root.search(searchParts, 0)
	if search != nil {
		parts := parsePattern(search.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return search, params
	}
	return nil, nil
}
func (r *router) handler(ctx *Context) {
	n, params := r.getRouter(ctx.Method, ctx.Path)
	if n != nil {
		ctx.Params = params
		key := ctx.Method + "-" + n.pattern
		r.handlers[key](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 not found:%s\n", ctx.Path)
	}
}
func (r *router) handlerMiddleware(ctx *Context) {
	n, params := r.getRouter(ctx.Method, ctx.Path)
	if n != nil {
		ctx.Params = params
		key := ctx.Method + "-" + n.pattern
		ctx.handlers = append(ctx.handlers, r.handlers[key])
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 not found:%s\n", ctx.Path)
		})
	}
	ctx.Next()
}
