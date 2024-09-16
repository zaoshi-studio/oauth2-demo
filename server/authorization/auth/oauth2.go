package auth

// generateCode creates a pseudo-random authorization code
func generateCode() string {
	return fmt.Sprintf("code-%d", time.Now().UnixNano())
}

// OAuth2AuthorizeHandler handles the /authorize endpoint
func OAuth2AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		clientID = r.FormValue("client_id")
	}
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = r.FormValue("redirect_uri")
	}
	responseType := r.URL.Query().Get("response_type")
	if responseType == "" {
		responseType = r.FormValue("response_type")
	}
	state := r.URL.Query().Get("state")
	if state == "" {
		state = r.FormValue("state")
	}

	// Validate response_type
	if responseType != "code" {
		http.Error(w, "Unsupported response_type", http.StatusBadRequest)
		return
	}

	// Validate client
	if _, err := db.GetAppInfoByID(clientID); err != nil {
		http.Error(w, "Invalid client_id", http.StatusUnauthorized)
		return
	}

	// Show login form
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "login.html") // Serve a simple login HTML form
		return
	}

	// Handle POST login
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		dbUser, err := db.GetUserByUsernameAndPassword(username, password)

		if err != nil || dbUser == nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate authorization code
		code := generateCode()

		// Store the code with associated user
		_ = redis.SetAuthorizationCode(code, dbUser.Username)

		// Redirect back to client with code and state
		redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// OAuth2TokenHandler handles the /token endpoint
func OAuth2TokenHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	if grantType != "authorization_code" {
		http.Error(w, "Unsupported grant_type", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	// Validate client
	if !helper.ValidateAppInfo(clientID, clientSecret) {
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	// Validate authorization code
	username, err := redis.GetAuthorizationCode(code)
	if err != nil {
		http.Error(w, "Invalid authorization code", http.StatusBadRequest)
		return
	}

	// Optionally, delete the code to prevent reuse
	redis.DelAuthorizationCode(code)

	// Create JWT access token
	token, err := helper.CreateJWT(username)
	if err != nil {
		http.Error(w, "Could not create access token", http.StatusInternalServerError)
		return
	}

	// Respond with token
	resp := struct {
		model.Token
	}{
		model.Token{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   3600, // 1 hour
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
