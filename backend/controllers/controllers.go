package controllers

import (
	"encoding/json"
	"log"

	"sugarweb.dev/web/backend/lib"
	"sugarweb.dev/web/backend/sugar"
)

func RootHandler(h *sugar.Handler) {
	data := struct {
		Content string `json:"content"`
		Success bool `json:"success"`
	} {
		Content: "created user",
		Success: true,
	}

	h.Response.JSON(data)
}

func SignupHandler(h *sugar.Handler) {
	var user lib.User
	if err := json.Unmarshal(h.Body, &user); err != nil {
		log.Fatal("Error")
	}

	log.Println(user.Email)
}