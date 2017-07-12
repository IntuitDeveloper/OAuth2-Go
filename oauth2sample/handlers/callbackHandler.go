package handlers

import (
	"log"
	"net/http"

	"oauth2sample/cache"
)

/*
 *  This is the redirect handler you configure in your app on developer.intuit.com
 *  The Authorization code has a short lifetime.
 *  Hence unless a user action is quick and mandatory, proceed to exchange the Authorization Code for
 *  BearerToken
 */
func CallBackFromOAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering CallBackFromOAuth ")
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	realmId := r.URL.Query().Get("realmId")

	cache.AddToCache("realmId", realmId)
	csrf := cache.GetFromCache("csrf")

	//check whether the state returned in the redirect is same as the csrf sent
	if state == csrf {

		// retrive bearer token using code
		bearerTokenResponse, err := RetrieveBearerToken(code)
		if err != nil {
			log.Fatalln(err)
		}
		/*
		 * add token to cache
		 * In real usecase, this is where tokens would have to be persisted (to a SQL DB, for example).
		 * Update your Datastore here with user's AccessToken and RefreshToken along with the realmId
		 */
		cache.AddToCache("access_token", bearerTokenResponse.AccessToken)
		cache.AddToCache("refresh_token", bearerTokenResponse.RefreshToken)

		/*
		 * However, in case of OpenIdConnect, when you request OpenIdScopes during authorization,
		 * you will also receive IDToken from Intuit. You first need to validate that the IDToken actually came from Intuit.
		 */
		idToken := bearerTokenResponse.IdToken
		if idToken != "" {
			//validate id token
			if ValidateIDToken(idToken) {
				// get userinfo
				GetUserInfo(w, r, bearerTokenResponse.AccessToken)
			}
		}
		log.Println("Exiting CallBackFromOAuth ")
		http.Redirect(w, r, "/connected/", http.StatusFound)
	}
}
