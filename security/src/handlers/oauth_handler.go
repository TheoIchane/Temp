package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/src/data"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)


type OAUTH_RESP struct {
	Type string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	Access_Token string `json:"access_token"`
	Error string `json:"error"`
}

type GoogleClient struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Verified_email bool `json:"verified_email"`
	Picture string `json:"picture"`
}

type GithubClient struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Avatar string `json:"avatar_url"`
}

type OAUTH_CONFIG struct {
	Client_ID string
	Client_Secret string
	Scope []string
	URL_Log string
	URL_Reg string
	URL string
	Endpoint Endpoint
}

type TokenReq struct {
	Client_ID string
	Client_Secret string
	GrantType string
}

type Endpoint struct {
	AuthURL string
	TokenURL string
	DeviceAuthURL string
	AuthStyle int
}

var Oauth_user *data.User
var Creds *data.Credentials

func OauthHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "google") {
		GoogleHandler(w, r )
		return
	} else if strings.Contains(r.URL.Path,"github"){
		GithubHandler(w,r)
		return
	}
}

func (c *OAUTH_CONFIG) GetEnv(app string) {
	if app == "google" {
		c.Client_ID = os.Getenv("GoogleID")
		c.Client_Secret = os.Getenv("GoogleSecret")
		c.URL = os.Getenv("GoogleURL")
		return
	}
	if app == "github" {
		c.Client_ID =  os.Getenv("GithubID")
		c.Client_Secret = os.Getenv("GithubSecret")
		c.URL = os.Getenv("GithubURL")
		return
	}
}

func (c OAUTH_CONFIG) ConfURL_Login(state string) string {
	buf := bytes.Buffer{}
	buf.WriteString(c.Endpoint.AuthURL+"?")
	values := url.Values{
		"response_type": {"code"},
		"client_id":     {c.Client_ID},
		"redirect_uri": {c.URL},
		"scope": {strings.Join(c.Scope," ")},
		"state": {"login"},
	}
	buf.WriteString(values.Encode())
	return buf.String()
}

func (c OAUTH_CONFIG) Token(code string,state string) string {
	fmt.Println(c)
	values := url.Values{
		"code": {code},
		"client_id": {c.Client_ID},
		"redirect_uri": {c.URL},
		"client_secret": {c.Client_Secret},
		"grant_type": {"authorization_code"},
	}
	buf:= bytes.NewBuffer([]byte(values.Encode()))
	Req,_ := http.NewRequest("POST",c.Endpoint.TokenURL,buf)
	Req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	Req.SetBasicAuth(url.QueryEscape(c.Client_ID), url.QueryEscape(c.Client_Secret))
	httpClient := &http.Client{}
	res, err := httpClient.Do(Req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	responseBuffer, _ := io.ReadAll(res.Body)
	responseJSON := string(responseBuffer)
	fmt.Println(responseJSON)
	return responseJSON
}

func (c OAUTH_CONFIG) ConfURL_Register(state string) string {
	buf := bytes.Buffer{}
	buf.WriteString(c.Endpoint.AuthURL+"?")
	values := url.Values{
		"response_type": {"code"},
		"client_id":     {c.Client_ID},
		"redirect_uri": {c.URL},
		"scope": {strings.Join(c.Scope," ")},
		"state": {state},
	}
	buf.WriteString(values.Encode())
	return buf.String()
}

func Client_Google(s string) *GoogleClient {
	resp := OAUTH_RESP{}
	_ = json.Unmarshal([]byte(s),&resp)
	if resp.Error != "" {
		fmt.Println(resp.Error)
		return nil
	}
	client := &http.Client{}
	res, err := client.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s",resp.Access_Token))
	if err != nil {
		fmt.Println(err)
		return &GoogleClient{}
	}
	buf , err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &GoogleClient{}
	}
	responseJson := string(buf)
	var google_client GoogleClient
	_ = json.Unmarshal([]byte(responseJson),&google_client)
	return &google_client
}

func Client_Github(s string) *GithubClient{
	token := strings.Split(strings.Split(s,"&")[0],"=")[1]
	req, err := http.NewRequest("GET","https://api.github.com/user",nil)
	val := fmt.Sprintf("token %s",token)
	req.Header.Set("Authorization",val)
	res, _ := http.DefaultClient.Do(req)
	buf , err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	responseJson := string(buf)
	fmt.Println(responseJson)
	var github_client GithubClient
	_ = json.Unmarshal([]byte(responseJson),&github_client)
	return &github_client
}
