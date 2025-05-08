package main

import (
	"fmt"

	"github.com/sam77il/sugar/db"
	"github.com/sam77il/sugar/middleware"
	"github.com/sam77il/sugar/sugar"
)

func main() {
	app := sugar.New(sugar.Config{
		Postgres: true,
	}, []sugar.Middleware{
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
	app.Post("/olala", olalaHandler)

	app.Listen(":8080")
}

func rootHandler(ctx *sugar.Context) {
	accounts := db.GetAllAccounts(ctx.DB, ctx.Request.Context())
	data := struct{
		Accounts []db.Account
	}{
		Accounts: accounts,
	}

	ctx.Page(data, "components/header", "components/footer", "pages/root")
}

func aHandler(ctx *sugar.Context) {
	ctx.HTML("<alo>ALO</alo>")
}

func olalaHandler(ctx *sugar.Context) {
	fmt.Println(ctx.Form())
	data := struct{
		Name string `json:"name"`
		Success bool `json:"success"`
	}{
		Name: ctx.Form().Get("name"),
		Success: true,
	}
	ctx.JSON(data)
}