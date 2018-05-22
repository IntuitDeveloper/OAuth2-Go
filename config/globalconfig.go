package config

import (
	"encoding/json"
	"log"
	"os"
)

//declare global variable and read from config
var OAuthConfig Config = LoadConfiguration("config.json")

//structure for configs
type Config struct {
	Port                        string `json:"port"`
	ClientId                    string `json:"clientId"`
	ClientSecret                string `json:"clientSecret"`
	RedirectUri                 string `json:"redirectUri"`
	IntuitAccountingAPIHost     string `json:"intuitAccountingAPIHost"`
	DiscoveryAPIHost            string `json:"discoveryAPIHost"`
	C2QBScope                   string `json:"c2qbScope"`
	SIWIScope                   string `json:"siwiScope"`
	GetAppNowScope              string `json:"getAppNowScope"`
	IntuitAuthorizationEndpoint string
}

/*
 * Method to load properties from config.json
 */
func LoadConfiguration(file string) Config {
	var conf Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&conf)

	return conf
}
