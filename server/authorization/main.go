package main

import (
	"fmt"
	"log"
	"net/http"
	"oauth2-demo/server/authorization/auth"
)

func main() {
	// Serve static login HTML
	http.HandleFunc("/auth/oauth2/authorize", auth.OAuth2AuthorizeHandler)
	http.HandleFunc("/auth/oauth2/token", auth.OAuth2TokenHandler)

	fmt.Println("Authorization Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
