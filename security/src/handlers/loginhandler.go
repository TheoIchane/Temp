package handlers

import (
	"encoding/base64"
	"fmt"
	"forum/src/data"
	"forum/src/middleware"
	"net/http"
	"time"
)

// LoginHandler gère la connexion des utilisateurs
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// If the user is already logged in, redirect to index
	userID, err := data.GetCurrentUserID(r)
	if err == nil && len(userID) != 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to process form", http.StatusBadRequest)
			return
		}
	}
	err_str := r.FormValue("error")
	Email := r.FormValue("email")
	Password := r.FormValue("password")
	p := &Page{Title: "Login", Data: make(map[string]interface{})}

	if Email == "" || Password == "" {
		if err_str == "login" {
			p.Error = "Ce compte n'existe pas, veuillez en créer un."
		}
		GITHUB_CONFIG.GetEnv("github")
		GOOGLE_CONFIG.GetEnv("google")
		p.Data["Google"] = GOOGLE_CONFIG.ConfURL_Login("200")
		p.Data["Github"] = GITHUB_CONFIG.ConfURL_Login("200")
		RenderTemplate(w, "login.html", p)
		return
	}
	if Email != "" {
		user, err := data.GetUserByEmail(Email)
		if err != nil || !middleware.ComparePassword(Password, user) {
			p.Error = "Invalid email or password"
			RenderTemplate(w, "login.html", p)
			return
		}

		Cookie(w, user)

		// Redirect to the index page after successful login
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	RenderTemplate(w, "login.html", p)
}

func Cookie(w http.ResponseWriter, user *data.User) {
	sessionID := data.GenerateUUID()

	session := data.Session{
		SessionID:  sessionID,
		UserUUID:   user.ID,
		Expires_at: time.Now().Add(1 * time.Hour),
	}

	err := data.InsertSession(&session)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(base64.StdEncoding.EncodeToString(sessionID))

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    base64.StdEncoding.EncodeToString(sessionID),
		Path:     "/",
		HttpOnly: true,
		Expires:  session.Expires_at,
		MaxAge:   3600,
	})
}
