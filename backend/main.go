package main

import (
	"sugarweb.dev/web/backend/controllers"
	"sugarweb.dev/web/backend/sugar"
)

func main() {
	app := sugar.New(sugar.Config{
		AppName: "MySite",
		Logs: true,
	})
	
	app.Get("/", controllers.RootHandler)
	app.Post("/auth/signup", controllers.SignupHandler)
	app.Listen(":7070")
}