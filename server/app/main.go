package main

import (
	"fmt"
	"log"
	"net/http"
	"oauth2-demo/server/app/auth"
	"oauth2-demo/server/app/handler"
)

func main() {
	http.HandleFunc("/resource", handler.ResourceHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/auth/oauth2/callback", auth.OAuth2CallbackHandler)

	fmt.Println("Resource Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
