package model

import (
	"encoding/json"
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
	json, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	return string(json)
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
