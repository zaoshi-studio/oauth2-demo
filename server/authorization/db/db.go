package db

import "oauth2-demo/pkg/model"

var apps = map[string]*model.AppInfo{
	"oauth2-demo-app": {
		ClientID:     "oauth2-demo-app",
		ClientSecret: "oauth2-demo-app-secret",
	},
}

var users = map[string]*model.User{
	"user1": {
		ID:       1,
		Username: "user1",
		Password: "password1",
	},
}

func GetAppInfoByID(clientID string) (*model.AppInfo, error) {
	return apps[clientID], nil
}

func GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	return users[username], nil
}
