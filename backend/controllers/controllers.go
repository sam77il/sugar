package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"sugarweb.dev/web/backend/lib"
	"sugarweb.dev/web/backend/sugar"
)

func RootHandler(h *sugar.Controller) {
	data := struct {
		Content string `json:"content"`
		Success bool `json:"success"`
	} {
		Content: "created user",
		Success: true,
	}

	h.Response.JSON(data)
}

func RootHandler2(h *sugar.Controller) {
	body := string(h.Request.Body)
	log.Println(body)
}

func SignupHandler(h *sugar.Controller) {
	var user lib.User
	if err := json.Unmarshal(h.Body, &user); err != nil {
		log.Fatal("Error")
	}

	log.Println(user.Email)
}

func LoginHandler(h *sugar.Controller) {
	var user lib.User
	if err := json.Unmarshal(h.Request.Body, &user); err != nil {
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

	cookie := http.Cookie{
		Name: "jwt",
		Value: token,
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	h.Response.SetCookie(&cookie)

	h.Response.JSON(map[string]bool{
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
