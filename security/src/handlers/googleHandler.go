package handlers

import (
	"fmt"
	"forum/src/data"
	"forum/src/middleware"
	"net/http"
)

var GOOGLE_CONFIG = &OAUTH_CONFIG{
	Scope:    []string{"email"},
	Endpoint: GoogleEndpoint,
}
var GoogleEndpoint = Endpoint{
	AuthURL:       "https://accounts.google.com/o/oauth2/auth",
	TokenURL:      "https://oauth2.googleapis.com/token",
	DeviceAuthURL: "https://oauth2.googleapis.com/device/code",
	AuthStyle:     0,
}

func GoogleHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	client := &GoogleClient{}
	var cred *data.Credentials
	var err error
	GOOGLE_CONFIG.GetEnv("google")
	if code != "" {
		resp := GOOGLE_CONFIG.Token(code)
		client = Client_Google(resp)
		if client == nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
	switch state {
	case "login":
		cred, err = data.GetCredentialsByID([]byte(client.ID))
		if err != nil {
			http.Redirect(w, r, "/login?error=login", http.StatusSeeOther)
			return
		}
		Oauth_user = cred.User
	case "register":
		fmt.Println(client.ID)
		if data.IsCredExist(middleware.Encrypt(client.ID)) {
			RenderTemplate(w, "register.html", Page{
				Title: "Register",
				Error: "Ce compte existe, veuillez vous connecter. ",
			})
			return
		}
		Oauth_user = &data.User{
			ID:       data.GenerateUUID(),
			Password: middleware.Encrypt(client.ID),
		}
		Creds = &data.Credentials{
			ID:         middleware.Encrypt(client.ID),
			Credential: "google",
		}
		http.Redirect(w, r, "/registeroauth", http.StatusSeeOther)
		return
	default:
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	Cookie(w, Oauth_user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterOauth(w http.ResponseWriter, r *http.Request) {
	p := Page{
		Title: "Login",
		Data: map[string]interface{}{
			"Oauth": true,
		},
	}
	if r.FormValue("username") != "" {
		Oauth_user.Username = r.FormValue("username")
		err := data.InsertUser(Oauth_user)
		if err == nil {
			Creds.User = Oauth_user
			err = Creds.Insert()
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
			Cookie(w, Oauth_user)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			p.Error = "Username already used"
		}
	}
	RenderTemplate(w, "oauth.html", p)
}
