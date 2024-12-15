package handlers

import (
	"fmt"
	"forum/src/data"
	"net/http"
	"strings"
)

// ProfileHandler gère
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	userID, err := data.GetCurrentUserID(r)
	username := r.URL.RawQuery
	profile := &data.User{}
	user, err := data.GetUserByUUID(userID)
	if username == "" || username == strings.ToLower(user.Username) {
		profile = user
	} else {
		profile, err = data.GetUserByUsername(username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	}
	datas := map[string]interface{}{
		"Liked_Posts": make([]map[string]interface{}, 0),
		"Posts":       make(map[string]interface{}, 0),
	}
	datas["Liked_Posts"] = data.GetLikedPosts(profile)
	datas["Posts"] = data.GetUserPosts(profile)
	datas["User"] = profile
	p := &Page{
		Title: "User Profile",
		User:  user,
		Data:  datas,
	}
	RenderTemplate(w, "profile.html", p)
}
