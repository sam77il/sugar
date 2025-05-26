package controllers

import (
	"encoding/json"
	"log"

	"sugarweb.dev/web/backend/lib"
	"sugarweb.dev/web/backend/sugar"
)

func RootHandler(crr *sugar.Controller) {
	data := struct {
		Content string `json:"content"`
		Success bool `json:"success"`
	} {
		Content: "created user",
		Success: true,
	}

	crr.Response.JSON(data)
}

func RootHandler2(crr *sugar.Controller) {
	body := string(crr.Request.Body)
	log.Println(body)
}

func SignupHandler(crr *sugar.Controller) {
	var user lib.User
	if err := json.Unmarshal(crr.Body, &user); err != nil {
		log.Fatal("Error")
	}

	log.Println(user.Email)
}

func LoginHandler(crr *sugar.Controller) {
	var user lib.User
	if err := json.Unmarshal(crr.Request.Body, &user); err != nil {
		log.Println(err.Error())
	}

	if user.Email != "admin@email.com" || user.Password != "abc123"  {
		log.Println("not admin")
		return
	}

	token, err := lib.GenereateJWT(user.Email, true)
	if err != nil {
		log.Println("Generating", err.Error())
		return
	}

	cookie := sugar.Cookie{
		Name: "jwt",
		Value: token,
		Path: "/",
		HttpOnly: true,
		SameSite: 3,
	}

	crr.Response.SetCookie(&cookie)

	crr.Response.JSON(map[string]bool{
		"success": true,
	})
}

func ProtectedHandler(h *sugar.Controller) {
	cookie, err := h.Request.Cookie("jwt")
	if err != nil {
		log.Println("Kein jwt token")
		return
	}

	claims, err := lib.ParseJWT(cookie.Value)
	if err != nil {
		log.Println("Token invalid")
		return
	}

	var user lib.User

	user.Email = (*claims)["email"].(string)
	user.IsAdmin = (*claims)["isAdmin"].(bool)

	log.Printf("%+v", user)
}
