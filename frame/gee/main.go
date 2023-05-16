package main

import (
	"fmt"
	"gee/gee"
	"log"
	"net/http"
	"time"
)

type Student struct {
	Name string
	Age  int
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
func onlyForV2() gee.HandlerFunc {
	return func(ctx *gee.Context) {
		now := time.Now()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(now))
	}
}
func main() {
	//r := gee.New()
	//################version 1
	//r.GET("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	//})
	//r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
	//	for k, v := range req.Header {
	//		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	//	}
	//})
	// ########### version 2
	//r.GET("/hello", func(ctx *gee.Context) {
	//	ctx.HTML(http.StatusOK, "<h1>HELLO GEE</h1>")
	//	ctx.String(http.StatusOK, "hello %s,you're at %s\n", ctx.Query("name"), ctx.Path)
	//})
	//r.POST("/login", func(ctx *gee.Context) {
	//	ctx.JSON(http.StatusOK, gee.H{
	//		"username": ctx.PostFrom("username"),
	//		"password": ctx.PostFrom("password"),
	//	})
	//})
	// ######version 3 实现版本
	//r.GET("/", func(ctx *gee.Context) {
	//	ctx.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//r.GET("/hello", func(ctx *gee.Context) {
	//	// expect   /hello?name=xuxinwen
	//	ctx.String(http.StatusOK, "hello %s,you're at %s\n", ctx.Query("name"), ctx.Path)
	//})
	//r.GET("hello/:name", func(ctx *gee.Context) {
	//	//except /hello/xuxinwen
	//	ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	//})
	//r.GET("/assets/*filepath", func(ctx *gee.Context) {
	//	ctx.JSON(http.StatusOK, gee.H{"filepath": ctx.Param("filepath")})
	//})
	//############分组实现 version4
	//r1 := gee.NewEngineGroup()
	//r.GET("/index", func(ctx *gee.Context) {
	//	ctx.HTML(http.StatusOK, "<h1>Index Page</h1>")
	//})
	//party := r1.Group("/v1")
	//{
	//	party.GET("/", func(ctx *gee.Context) {
	//		ctx.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//	})
	//	party.GET("/hello", func(ctx *gee.Context) {
	//		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	//	})
	//}
	//party1 := r1.Group("/v2")
	//{
	//	party1.GET("/hello/:name", func(ctx *gee.Context) {
	//		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	//	})
	//	party1.GET("/login", func(ctx *gee.Context) {
	//		ctx.JSON(http.StatusOK, gee.H{
	//			"username": ctx.PostFrom("username"),
	//			"password": ctx.PostFrom("password"),
	//		})
	//	})
	//}
	//r1.Run(":9999")
	//############### version5
	//r := gee.NewEngineGroup()
	//r.Use(gee.Logger()) //global middleware
	//r.GET("/", func(ctx *gee.Context) {
	//	ctx.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//v2 := r.Group("/v2")
	//v2.Use(onlyForV2()) // v2 group middleware
	//{
	//	v2.GET("/hello/:name", func(ctx *gee.Context) {
	//		// expect /hello/geektutu
	//		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	//	})
	//}
	//r.Run(":9999")
	//############template version6
	//r := gee.NewEngineGroup()
	//r.Use(gee.Logger())
	//r.SetFuncMap(template.FuncMap{"FormatAsDate": FormatAsDate})
	//r.LoadHTMLGlob("templates/*")
	//r.Static("/assets", "./static")
	//stu1 := &Student{Name: "xxw", Age: 26}
	//stu2 := &Student{Name: "zzy", Age: 46}
	//r.GET("/", func(ctx *gee.Context) {
	//	ctx.HTMLTemplate(http.StatusOK, "css.tmpl", nil)
	//})
	//r.GET("/students", func(ctx *gee.Context) {
	//	ctx.HTMLTemplate(http.StatusOK, "arr.tmpl", gee.H{
	//		"title":  "gee",
	//		"stuArr": [2]*Student{stu1, stu2},
	//	})
	//})
	//r.GET("/date", func(ctx *gee.Context) {
	//	ctx.HTMLTemplate(http.StatusOK, "custom_func.tmpl", gee.H{
	//		"title": "gee",
	//		"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
	//	})
	//})
	//r.Run(":9999")
	//####### 错误处理
	r := gee.NewEngineGroup()
	middleware := make([]gee.HandlerFunc, 0)
	middleware = append(middleware, gee.Recovery(), gee.Logger())
	r.Use(middleware...)
	r.GET("/", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello Geektutu\n")
	})
	r.GET("/panic", func(ctx *gee.Context) {
		names := []string{"xxw"}
		ctx.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}
