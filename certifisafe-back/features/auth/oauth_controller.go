package auth

import (
	"certifisafe-back/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type OauthController struct {
	oauthService OauthService
	oauthConfig  *oauth2.Config
	oauthUrlAPI  string
	clientURL    string
}

func NewOauthController(oauthService OauthService) *OauthController {
	config := utils.Config()
	return &OauthController{
		oauthService: oauthService,
		oauthConfig: &oauth2.Config{
			ClientID:     config["oauth-client-id"],
			ClientSecret: config["oauth-client-secret"],
			RedirectURL:  config["oauth-redirect-url"],
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"}, // Add required scopes
			Endpoint:     google.Endpoint,
		},
		oauthUrlAPI: config["oauth-api-url"],
		clientURL:   config["client-url"],
	}
}

func (controller *OauthController) Oauth(w http.ResponseWriter, r *http.Request) {
	oauthState := controller.generateStateOauthCookie(w)
	u := controller.oauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (controller *OauthController) OauthCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := controller.getUserDataFromOauthService(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := controller.oauthService.AuthenticateUser(data)
	query := url.Values{}
	query.Add("token", token)
	controller.clientURL += "?" + query.Encode()
	http.Redirect(w, r, controller.clientURL, http.StatusMovedPermanently)
}

func (controller *OauthController) generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(time.Minute * 60)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (controller *OauthController) getUserDataFromOauthService(code string) ([]byte, error) {
	token, err := controller.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(controller.oauthUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
