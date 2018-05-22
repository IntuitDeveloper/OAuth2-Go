package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/IntuitDeveloper/OAuth2-Go/cache"
	"github.com/IntuitDeveloper/OAuth2-Go/config"
)

/*
 * Method to retrive access token (bearer token)
 */
func RetrieveBearerToken(code string) (*BearerTokenResponse, error) {
	log.Println("Entering RetrieveBearerToken ")
	client := &http.Client{}
	data := url.Values{}
	//set parameters
	data.Set("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", config.OAuthConfig.RedirectUri)

	tokenEndpoint := cache.GetFromCache("token_endpoint")
	request, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	//set headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth())

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	bearerTokenResponse, err := getBearerTokenResponse([]byte(body))
	log.Println("Exiting RetrieveBearerToken ")
	return bearerTokenResponse, err
}

type BearerTokenResponse struct {
	RefreshToken           string `json:"refresh_token"`
	AccessToken            string `json:"access_token"`
	TokenType              string `json:"token_type"`
	IdToken                string `json:"id_token"`
	ExpiresIn              int64  `json:"expires_in"`
	XRefreshTokenExpiresIn int64  `json:"x_refresh_token_expires_in"`
}

func getBearerTokenResponse(body []byte) (*BearerTokenResponse, error) {
	var s = new(BearerTokenResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln("error getting BearerTokenResponse:", err)
	}
	return s, err
}

func basicAuth() string {
	auth := config.OAuthConfig.ClientId + ":" + config.OAuthConfig.ClientSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
