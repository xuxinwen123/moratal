package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type H map[string]interface{}
type Context struct {
	W   http.ResponseWriter
	Req *http.Request
	//Req info
	Path   string
	Method string
	Params map[string]string
	//Resp
	StatusCode int
	//middleware
	//中间件不仅作用在处理流程前，也可以作用在处理流程后 ,所以可以放在context里面
	handlers []HandlerFunc
	index    int //记录当前执行到第几个中间件，当在中间件中调用Next方法的时候，控制权移交给下一个中间件，
	// 直到调用到最后一个中间件，然后再从后往前，调用每个中间件在Next方法之后定义的那部分。
	engine *EngineGroup // engine pointer
}

func (ctx *Context) Param(key string) string {
	s := ctx.Params[key]
	return s
}
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		W: w, Req: req, Path: req.URL.Path, Method: req.Method, index: -1}
}

// Next 中间件可等待用户自己定义的 Handler处理结束后，再做一些额外的操作
func (ctx *Context) Next() {
	ctx.index++
	//不是所有的handler都会调用 Next()。
	//手工调用 Next()，一般用于在请求前后各实现一些行为。如果中间件只作用于请求前，可以省略调用Next()，此种写法可以兼容 不调用Next的写法
	//并且使用c.index<s;c.index++ 也能保证中间件只执行一次
	//当中间件不调用 next函数时,通过此循环保证中间件执行顺序
	//当中间件调用next 函数时,且当中间件执行完毕,对应c.index也已经到达指定index
	//前面存在的for循环因为是通过c.index<s 做循环判断,则不会重复执行已经执行过的中间件
	s := len(ctx.handlers)
	for ; ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}

}
func (ctx *Context) PostFrom(key string) string {
	return ctx.Req.FormValue(key)
}
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.W.WriteHeader(code)
}
func (ctx *Context) SetHeader(key, value string) {
	ctx.W.Header().Set(key, value)
}
func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.W.Write([]byte(fmt.Sprintf(format, values...)))
}
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Context-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.W)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.W, err.Error(), 500)
	}
}
func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.W.Write(data)
}
func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("Content_type", "text/html")
	ctx.Status(code)
	ctx.W.Write([]byte(html))
}
func (ctx *Context) HTMLTemplate(code int, name string, data interface{}) {
	ctx.SetHeader("Content_type", "text/html")
	ctx.Status(code)
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.W, name, data); err != nil {
		log.Fatal(err.Error())
	}
}
