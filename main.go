package main

import (
	"bytes"
	"context"
	"html/template"
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

	app.Listen(8080)
}

type User struct {
	Email string
	Password string
}

func rootHandler(ctx *sugar.Context) {
	tmpl := template.Must(template.ParseFiles(
	    "components/header.sugar",
	    "components/footer.sugar",
	    "pages/root.sugar",
	))

	rows, err := ctx.DB.Query(ctx.Request.Context(), "SELECT email, password FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
		
	var accounts []User
	for rows.Next() {
		var acc User
		if err := rows.Scan(&acc.Email, &acc.Password); err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, acc)
	}

	data := struct{
		Users []User
	}{
		Users: accounts,
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "content", data)
	if err != nil {
		log.Fatal(err)
	}
	ctx.HTML(buf.String())
}

func aHandler(ctx *sugar.Context) {
	body, header, err := ctx.Get("http://localhost:8080/b")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(body, header.Get("Content-Type"))
}