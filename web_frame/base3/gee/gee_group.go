package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
	//支持模板渲染
	htmlTemplates *template.Template //for html render 将所有的模板加载进内存
	funcMap       template.FuncMap   //for html render 所有的自定义模板渲染函数
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

// create static handler
// 用户可以将磁盘上的某个文件夹root映射到路由relativePath
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		//check if file exist and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.W, ctx.Req)
	}
}

// Static serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	//Register Get handlers
	group.GET(urlPattern, handler)
}
func (e *EngineGroup) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}
func (e *EngineGroup) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
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
	//使用中间件需要增加的部分
	e.engine = e
	e.router.handlerMiddleware(ctx)
	//给engine 赋值
	//分组的实现
	//ctx := newContext(w, r)
	//group.router.handler(ctx)
}
