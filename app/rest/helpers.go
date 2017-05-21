package rest

import "crypto/subtle"

func CheckSecretMatch(x, y string) bool {
	return subtle.ConstantTimeCompare([]byte(x), []byte(y)) == 1
}
