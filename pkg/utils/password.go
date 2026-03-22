package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw),12 /*bcrypt.DefaultCost*/)
	return string(hash), err
}

func CheckPassword(hash, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}