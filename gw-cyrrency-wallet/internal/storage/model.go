package storage

import "fmt"

type Account struct {
	username string
	password string
	email    string
}

func (a *Account) String() string {
	out := fmt.Sprintf("username = %v\t", a.username)
	out += fmt.Sprintf("password = %v\n", a.password)
	return out
}
