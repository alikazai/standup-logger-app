package utils

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func Password_verify(hashed string, raw_pwd []byte) bool {
	bytehash := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(bytehash, raw_pwd) // Compare hash
	if err != nil {
		log.Err(err).Msg("password not matched")
		return false
	}
	return true
}
