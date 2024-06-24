package auth

import (
	"crypto/sha512"
	"encoding/hex"
)

func hashPassword(password string, salt []byte) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, salt...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}

func matchPassword(hashedPassword, currPassword string, salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)
	return hashedPassword == currPasswordHash
}
