package main

import (
	"github.com/joho/godotenv"
	"sugarweb.dev/web/backend/controllers"
	"sugarweb.dev/web/backend/sugar"
)

func main() {
	godotenv.Load()
	app := sugar.New(sugar.Config{
		AppName: "MySite",
		Logs: true,
	})
	
	app.Get("/", controllers.RootHandler)
	app.Post("/:id", controllers.RootHandler2)
	app.Post("/auth/signup", controllers.SignupHandler)

	app.Listen(7070)
}