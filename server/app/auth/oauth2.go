package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"oauth2-demo/pkg/model"
	"oauth2-demo/server/app/config"
)

// OAuth2CallbackHandler Callback after receiving authorization code
func OAuth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")

	// Exchange the authorization code for a token
	resp, err := http.PostForm(config.OAuth2TokenURL,
		map[string][]string{
			"client_id": {
				config.OAuth2ClientID,
			},
			"client_secret": {
				config.OAuth2ClientSecret,
			},
			"grant_type": {
				config.OAuth2GrantType,
			},
			"code": {authCode},
		})
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the token response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var tokenResponse model.Token
	json.Unmarshal(body, &tokenResponse)

	// Display token
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
