package main

import (
	"log"

	"github.com/sam77il/dusch-wand/middleware"
	"github.com/sam77il/dusch-wand/sugar"
)

func main() {
	app := sugar.New([]sugar.Middleware{
		{
			Path: "/a",
			MiddlewareFunc: middleware.LogRoute,
		},
	})
	
	app.Get("/", rootHandler)
	app.Get("/a", aHandler)
	app.Get("/b", rootHandler)
	app.Get("/c", rootHandler)
	app.Get("/favicon.ico", func(ctx *sugar.Context) {
		ctx.NotFound()
	})

	app.Listen(8080)
}

func rootHandler(ctx *sugar.Context) {
	ctx.HTML("test")
}

func aHandler(ctx *sugar.Context) {
	body, header, err := ctx.Get("http://localhost:8080/b")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(body, header.Get("Content-Type"))
}