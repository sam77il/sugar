package main

import (
	"github.com/sam77il/sugar/db"
	"github.com/sam77il/sugar/middleware"
	"github.com/sam77il/sugar/sugar"
)

func main() {
	app := sugar.New(sugar.Config{
		Postgres: true,
	})

	app.Middleware("/", middleware.LogRoute)
	
	app.Get("/", rootHandler)

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