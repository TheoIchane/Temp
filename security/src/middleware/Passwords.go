package middleware

import (
	"forum/src/data"

	"golang.org/x/crypto/bcrypt"
)

//Encrypt a string by Hashing it
func Encrypt(str string) []byte {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), 15)
	if err != nil {
	}
	return hashed
}

func ComparePassword(password string, user *data.User) bool {
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return false
	}
	return true
}
