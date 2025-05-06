package main

import (
	"github.com/sam77il/dusch-wand/middleware"
	"github.com/sam77il/dusch-wand/sugar"
)

func main() {
	app := sugar.New([]sugar.Middleware{
		{
			Path: "*",
			MiddlewareFunc: middleware.LogRoute,
		},
	})
	
	app.Get("/", rootHandler)
	app.Get("/a", rootHandler)
	app.Get("/b", rootHandler)
	app.Get("/c", rootHandler)

	app.Listen(8080)
}

func rootHandler(ctx *sugar.Context) {
	ctx.HTML("test")
}