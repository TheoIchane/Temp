package handlers

import (
	"forum/src/data"
	"forum/src/middleware"
	"net/http"
	"strings"
)

// RegisterHandler gère l'enregistrement des nouveaux utilisateurs
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	userID, err := data.GetCurrentUserID(r)
	if err == nil && len(userID) > 0 {
		// User is logged in, redirect to index page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var errorMessages []string
	var user data.User

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			errorMessages = append(errorMessages, "Unable to process form")
		} else {
			// Validation des données du formulaire
			err_str := r.FormValue("error")
			email := r.FormValue("email")
			password := r.FormValue("password")
			username := r.FormValue("username")
			if err_str == "register" {
				errorMessages = append(errorMessages, "Ce compte existe dejà, veuillez vous connecter.")
			} else {
				// Vérification de la validité de l'email
				if !middleware.IsValidEmail(email) {
					errorMessages = append(errorMessages, "Invalid Email.")
				}

				UsernameErrors := middleware.IsValidUsername(username)
				if len(UsernameErrors) > 0 {
					errorMessages = append(errorMessages, UsernameErrors...)
				}

				// Validation du mot de passe
				passwordErrors := middleware.ValidatePassword(password)
				if len(passwordErrors) > 0 {
					errorMessages = append(errorMessages, passwordErrors...)
				}
			}
			// Si aucune erreur, créer un nouvel utilisateur
			if len(errorMessages) == 0 {
				var admin bool
				user = data.User{
					ID:       []byte(data.GenerateUUID()),
					Email:    email,
					Username: username,
					Password: middleware.Encrypt(password), // Encrypt password here with Password encrypt
					Avatar:   "",
					Admin:    admin,
				}
				err := data.InsertUser(&user)
				if err != nil {
					errorMessages = append(errorMessages, err.Error())
				} else {

					Cookie(w, &user)
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}
			}
		}
	}
	// Affichage des erreurs dans la page d'enregistrement
	p := &Page{
		Title: "Register",
		Error: strings.Join(errorMessages, "<br>"), // Séparer les erreurs par des sauts de ligne HTMLù
		Data:  make(map[string]interface{}),
	}
	GITHUB_CONFIG.GetEnv("github")
	GOOGLE_CONFIG.GetEnv("google")
	p.Data["Google"] = GOOGLE_CONFIG.ConfURL_Register()
	p.Data["Github"] = GITHUB_CONFIG.ConfURL_Register()
	RenderTemplate(w, "register.html", p)
}
