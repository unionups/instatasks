package test_helpers

import (
	"encoding/base64"
)

func AuthorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}
