package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sam77il/sugar/middleware"
	"github.com/sam77il/sugar/sugar"
)

func main() {
	godotenv.Load()
	db, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := sugar.New(db, []sugar.Middleware{
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
	app.Get("/olala", olalaHandler)

	app.Listen(8080)
}

type Account struct {
	Email string
	Password string
}

func rootHandler(ctx *sugar.Context) {
	var accounts []Account
	rows, err := ctx.DB.Query(ctx.Request.Context(), "SELECT email, password FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var account Account

		rows.Scan(&account.Email, &account.Password)
		accounts = append(accounts, account)
	}

	data := struct{
		Accounts []Account
	}{
		Accounts: accounts,
	}

	ctx.Page(data, "components/header", "components/footer", "pages/root")
}

func aHandler(ctx *sugar.Context) {
	ctx.HTML("<alo>ALO</alo>")
}

func olalaHandler(ctx *sugar.Context) {
	fmt.Println(ctx.Form().Get("name"))
}