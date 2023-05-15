package gee

import (
	"encoding/json"
	"fmt"
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
