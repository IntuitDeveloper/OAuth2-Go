package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"oauth2sample/cache"
	"oauth2sample/config"
)

/*
 * Sample QBO API call to get CompanyInfo using OAuth2 tokens
 */
func GetCompanyInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering GetCompanyInfo ")
	client := &http.Client{}

	//Ideally you would fetch the realmId and the accessToken from the data store based on the user account here.
	realmId := cache.GetFromCache("realmId")
	if realmId == "" {
		log.Println("No realm ID.  QBO calls only work if the accounting scope was passed!")
		fmt.Fprintf(w, "No realm ID.  QBO calls only work if the accounting scope was passed!")
	}
	request, err := http.NewRequest("GET", config.OAuthConfig.IntuitAccountingAPIHost+"/v3/company/"+realmId+"/companyinfo/"+realmId+"?minorversion=8", nil)
	if err != nil {
		log.Fatalln(err)
	}
	//set header
	request.Header.Set("accept", "application/json")
	accessToken := cache.GetFromCache("access_token")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	responseString := string(body)
	log.Println("Exiting GetCompanyInfo ")
	fmt.Fprintf(w, responseString)

}
