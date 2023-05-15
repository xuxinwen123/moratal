package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

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
	r := gee.NewEngineGroup()
	r.Use(gee.Logger()) //global middleware
	r.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(ctx *gee.Context) {
			// expect /hello/geektutu
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
	}
	r.Run(":9999")
}
