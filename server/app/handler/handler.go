package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"oauth2-demo/pkg/helper"
	"oauth2-demo/server/app/config"
	"strings"
)

// LoginHandler Simulated App Server - Redirect to Authorization Server
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate redirect to authorization server
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Server internal error", http.StatusBadRequest)
		return
	}

	channel := r.FormValue("channel")
	switch channel {
	case "oauth2":
		authURL := fmt.Sprintf("%s?response_type=%v&client_id=%v&redirect_uri=%v&username=player1&password=password123&state=xyz",
			config.OAuth2AuthURL,
			"code",
			config.OAuth2ClientID,
			config.OAuth2Callback,
		)
		http.Redirect(w, r, authURL, http.StatusPermanentRedirect)
	case "credential":

	default:
		http.Error(w, "Invalid channel", http.StatusBadRequest)
	}
}

// ResourceHandler handles requests to the /resource endpoint
func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Check if the Authorization header is in the correct format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	// Validate the JWT token
	token, err := helper.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Optionally, extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Invalid token subject", http.StatusUnauthorized)
		return
	}

	// Provide the protected resource
	response := struct {
		Data string `json:"data"`
	}{
		Data: fmt.Sprintf("Hello, %s! This is your protected game data.", userID),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
