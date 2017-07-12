package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"oauth2sample/cache"
)

/*
 * Call the revoke endpoint to revoke tokens
 */
func RevokeToken(w http.ResponseWriter, r *http.Request) {

	log.Println("Entering RevokeToken ")
	client := &http.Client{}
	data := url.Values{}

	//add parameters
	refreshToken := cache.GetFromCache("refresh_token")
	data.Add("token", refreshToken)

	revokeEndpoint := cache.GetFromCache("revocation_endpoint")
	request, err := http.NewRequest("POST", revokeEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	//set headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth())

	resp, err := client.Do(request)
	defer resp.Body.Close()

	responseString := map[string]string{"response": "Revoke successful"}
	responseData, _ := json.Marshal(responseString)
	log.Println("Exiting RevokeToken ")
	fmt.Fprintf(w, string(responseData))

}
