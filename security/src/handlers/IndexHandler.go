package handlers

import (
	"fmt"
	"forum/src/data"
	"net/http"
)

// IndexHandler affiche la page d'accueil
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		// No valid cookie, user is not logged in
		fmt.Println("User is not logged in")
		RenderIndex(w, r)
		return
	}

	isValid, err := data.IsSessionValid(cookie.Value)
	if err != nil {
		// Handle database or validation error
		fmt.Println("Error validating session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !isValid {
		// Invalid session, clear the cookie and redirect to login
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1, // Expire the cookie immediately
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	RenderIndex(w, r)
}

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	var user *data.User
	userID, err := data.GetCurrentUserID(r)
	if err == nil {
		user, err = data.GetUserByUUID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
		}
	}
	datas := map[string]interface{}{
		"Posts":       make([]map[string]interface{}, 0),
	}
	datas["Posts"] = data.GetLastPosts(user)
	p := &Page{
		Title: "Forum",
		User:  user,
		Data:  datas,
	}
	RenderTemplate(w, "index.html", p)
}
