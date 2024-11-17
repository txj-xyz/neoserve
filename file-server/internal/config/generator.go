package config

import (
	"crypto/rand"
	"encoding/hex"
)

/*
	Generates a random name that is 16 bytes

	@returns [string]
*/
func GenerateRandomName() (string) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil { return "" }
	r := hex.EncodeToString(randomBytes)
	return r
}