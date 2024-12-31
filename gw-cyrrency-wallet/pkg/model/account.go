package model

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Username string
	Password string
	Email    string
}

func (a *Account) String() string {
	out := "{\n"
	out += fmt.Sprintf("\tusername: %v\n", a.Username)
	out += fmt.Sprintf("\tpassword: %v\n", a.Password)
	out += fmt.Sprintf("\temail: %v\n", a.Email)
	out += "}"
	return out
}

func (a *Account) SetPassword() ([]byte, error) {
	salt := os.Getenv("SALT_DB")
	secret := a.Password + salt
	return bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
}

func (a *Account) CheckPassword(passwordCript string) bool {
	salt := os.Getenv("SALT_DB")
	secret := a.Password + salt
	return bcrypt.CompareHashAndPassword([]byte(passwordCript), []byte(secret)) == nil
}
