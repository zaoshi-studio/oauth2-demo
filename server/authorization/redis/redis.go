package redis

var authorizationCodes = map[string]string{}

func SetAuthorizationCode(code, username string) error {
	authorizationCodes[code] = username
	return nil
}

func GetAuthorizationCode(code string) (string, error) {
	return authorizationCodes[code], nil
}

func DelAuthorizationCode(code string) error {
	delete(authorizationCodes, code)
	return nil
}
