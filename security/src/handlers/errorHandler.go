package handlers

import "net/http"

// ErrorHandler gère les erreurs HTTP
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	p := &Page{Title: "Error", Error: http.StatusText(status)}
	RenderTemplate(w, "error.html", p)
}
