package main

import (
	"bytes"
	"html/template"
	"log"

	"github.com/sam77il/sugar/middleware"
	"github.com/sam77il/sugar/sugar"
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

type User struct {
	Email string
	Password string
}

func rootHandler(ctx *sugar.Context) {
	tmpl := template.Must(template.ParseFiles(
	    "layouts/root.sugar",
	    "components/header.sugar",
	    "components/footer.sugar",
	    "pages/root.sugar",
	))
		
	data := struct {
	    Head_Title  string
		Users []User
	}{
	    Head_Title: "My Page",
		Users: []User{
			{
				Email: "yavuzsamil.guengoer@gmail.com",
				Password: "test123",
			},
			{
				Email: "kingyavuzea7@gmail.com",
				Password: "test1234",
			},
			{
				Email: "samil_yavuz@web.de",
				Password: "test12345",
			},
			{
				Email: "revenmainacc@gmail.com",
				Password: "test123456",
			},
		},
	}
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "layout", data)
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