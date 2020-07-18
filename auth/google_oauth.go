package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	googleOauthConfig oauth2.Config
	httpClient        *http.Client
}

func New(_oauth oauth2.Config) Auth {
	var a Auth
	a.googleOauthConfig = _oauth
	a.googleOauthConfig.Scopes = []string{"https://www.googleapis.com/auth/userinfo.email"}
	a.googleOauthConfig.Endpoint = google.Endpoint
	a.httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	return a
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="

func (o Auth) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) string {

	// Create oauthState cookie
	oauthState := o.generateStateOauthCookie(w)

	/*
	   AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
	   validate that it matches the the state query parameter on your redirect callback.
	*/
	return o.googleOauthConfig.AuthCodeURL(oauthState)

}

func (o Auth) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) (string, error) {

	// Read oauthState from Cookie
	/*oauthState, err := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return err
	}*/

	token, err := o.getGoogleToken(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return "", err
	}

	return token, err
}

func (o Auth) generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	return state
}

func (o Auth) getGoogleToken(code string) (string, error) {
	// Use code to get token from Google.
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(&map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     o.googleOauthConfig.ClientID,
		"client_secret": o.googleOauthConfig.ClientSecret,
		"redirect_uri":  o.googleOauthConfig.RedirectURL,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, `https://www.googleapis.com/oauth2/v4/token`, &buffer)
	if err != nil {
		return "", err
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Token string `json:"id_token"`
	}

	err = jsoniter.ConfigFastest.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.Token, nil
}

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail string `json:"verified_email"`
}
