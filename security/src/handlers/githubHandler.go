package handlers

import (
	"forum/src/data"
	"forum/src/middleware"
	"net/http"
	"strconv"
)

var GITHUB_CONFIG = OAUTH_CONFIG{
	Endpoint: GithubEndpoint,
}

var GithubEndpoint = Endpoint{
	AuthURL:       "https://github.com/login/oauth/authorize",
	TokenURL:      "https://github.com/login/oauth/access_token",
	DeviceAuthURL: "https://github.com/login/device/code",
	AuthStyle:     0,
}

func GithubHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	client := &GithubClient{}
	if code != "" {
		res := GITHUB_CONFIG.Token(code)
		client = Client_Github(res)
	}
	Oauth_user, _ = data.GetUserByEmail(client.Email)
	if len(Oauth_user.ID) == 0 {
		Oauth_user.ID = data.GenerateUUID()
		Oauth_user.Email = client.Email
		Oauth_user.Password = middleware.Encrypt(strconv.Itoa(client.ID))
		RegisterOauth(w, r)
		return
	}

	Cookie(w, Oauth_user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
