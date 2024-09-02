package utills

import "encoding/base64"

func CreateBasicAuthToken(username, password string) string {
	credentials := username + ":" + password
	token := base64.StdEncoding.EncodeToString([]byte(credentials))
	return token
}
