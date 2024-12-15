package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookie("session_id"))
	// Clear the session cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire the cookie
		MaxAge:   -1,                             // Set MaxAge to -1 to delete the cookie
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	// Redirect the user to the login page after logging out
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
