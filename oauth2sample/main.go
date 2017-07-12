package main

import (
	"log"
	"net/http"
	"oauth2sample/config"
	"oauth2sample/handlers"
)

func main() {

	//call discovery
	handlers.CallDiscoveryAPI()

	//register static routes
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/connected/", http.StripPrefix("/connected/", http.FileServer(http.Dir("static/connected/"))))

	//register handler routes
	http.HandleFunc("/getCompanyInfo", handlers.GetCompanyInfo)
	http.HandleFunc("/refreshToken", handlers.RefreshToken)
	http.HandleFunc("/revokeToken", handlers.RevokeToken)
	http.HandleFunc("/connectToQuickbooks", handlers.ConnectToQuickbooks)
	http.HandleFunc("/signInWithIntuit", handlers.SignInWithIntuit)
	http.HandleFunc("/getAppNow", handlers.GetAppNow)
	http.HandleFunc("/oauth2redirect", handlers.CallBackFromOAuth)

	//log and start server
	log.Println("running server on ", config.OAuthConfig.Port)
	log.Fatal(http.ListenAndServe(config.OAuthConfig.Port, nil))

}
