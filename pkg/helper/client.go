package helper

import "oauth2-demo/server/authorization/db"

// ValidateAppInfo verifies client credentials
func ValidateAppInfo(clientID, clientSecret string) bool {

	dbData, err := db.GetAppInfoByID(clientID)
	if err != nil {
		return false
	}

	return dbData.ClientID == clientID && dbData.ClientSecret == clientSecret
}
