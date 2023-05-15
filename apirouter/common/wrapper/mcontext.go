package wrapper

import (
	"context"
	"sync"
)

// Context is our custom context.
// Let's implement a context which will give us access
type Context struct {
	context.Context
	//添加需要的东西 eg:username
}

var contextPool = sync.Pool{New: func() interface{} {
	return &Context{}
}}

func acquire(original context.Context) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Context = original

	//设置其它信息
	return ctx
}
