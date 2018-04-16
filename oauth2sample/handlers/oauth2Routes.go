package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/twinj/uuid"

	"oauth2sample/cache"
	"oauth2sample/config"
)

/*
 * Handler for connectToQuickbooks button
 */
func ConnectToQuickbooks(w http.ResponseWriter, r *http.Request) {
	log.Println("inside connectToQuickbooks ")
	http.Redirect(w, r, PrepareUrl(config.OAuthConfig.C2QBScope, GenerateCSRF()), http.StatusSeeOther)
}

/*
 * Handler for signInWithIntuit button
 */
func SignInWithIntuit(w http.ResponseWriter, r *http.Request) {
	log.Println("inside signInWithIntuit ")
	http.Redirect(w, r, PrepareUrl(config.OAuthConfig.SIWIScope, GenerateCSRF()), http.StatusSeeOther)
}

/*
 * Handler for getAppNow button
 */
func GetAppNow(w http.ResponseWriter, r *http.Request) {
	log.Println("inside getAppNow ")
	http.Redirect(w, r, PrepareUrl(config.OAuthConfig.GetAppNowScope, GenerateCSRF()), http.StatusSeeOther)
}

/*
 * Generates CSRF token
 */
func GenerateCSRF() string {
	csrf := uuid.NewV4().String()
	//add to cache since we need this in callback handler to validate the response
	cache.AddToCache("csrf", csrf)
	return csrf
}

/*
 * Prepares URL to call the OAuth2 authorization endpoint using Scope, CSRF and redirectURL that is supplied
 */
func PrepareUrl(scope string, csrf string) string {
	var Url *url.URL

	authorizationEndpoint := cache.GetFromCache("authorization_endpoint")
	Url, err := url.Parse(authorizationEndpoint)
	if err != nil {
		panic("error parsing url")
	}

	parameters := url.Values{}
	parameters.Add("client_id", config.OAuthConfig.ClientId)
	parameters.Add("response_type", "code")
	parameters.Add("scope", scope)
	parameters.Add("redirect_uri", config.OAuthConfig.RedirectUri)
	parameters.Add("state", csrf)
	Url.RawQuery = parameters.Encode()

	log.Printf("Encoded URL is %q\n", Url.String())
	return Url.String()

}
