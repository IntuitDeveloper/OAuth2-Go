package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/IntuitDeveloper/OAuth2-Go/cache"
	"github.com/IntuitDeveloper/OAuth2-Go/config"
)

/*
 *  Call discovery document and populate the cache
 */
func CallDiscoveryAPI() {
	log.Println("Entering CallDiscoveryAPI ")
	client := &http.Client{}
	request, err := http.NewRequest("GET", config.OAuthConfig.DiscoveryAPIHost, nil)
	if err != nil {
		log.Fatalln(err)
	}
	//set header
	request.Header.Set("accept", "application/json")

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	discoveryAPIResponse, err := getDiscoveryAPIResponse([]byte(body))

	//Add the urls to cache - in real app, these should be stored in database or config repository
	cache.AddToCache("authorization_endpoint", discoveryAPIResponse.AuthorizationEndpoint)
	cache.AddToCache("token_endpoint", discoveryAPIResponse.TokenEndpoint)
	cache.AddToCache("jwks_uri", discoveryAPIResponse.JwksUri)
	cache.AddToCache("revocation_endpoint", discoveryAPIResponse.RevocationEndpoint)
	cache.AddToCache("userinfo_endpoint", discoveryAPIResponse.UserinfoEndpoint)
	cache.AddToCache("issuer", discoveryAPIResponse.Issuer)

	log.Println("Exiting CallDiscoveryAPI ")
}

type DiscoveryAPIResponse struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
	RevocationEndpoint    string `json:"revocation_endpoint"`
	JwksUri               string `json:"jwks_uri"`
}

func getDiscoveryAPIResponse(body []byte) (*DiscoveryAPIResponse, error) {
	var s = new(DiscoveryAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln("error getting DiscoveryAPIResponse:", err)
	}
	return s, err
}
