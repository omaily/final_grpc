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

func (st *Account) CreateAccount(name, pass string) (*Account, bool) {
	user := Account{
		username: name,
		password: pass,
	}
	return &user, true
}

func (st *Instance) FindAccount(name, pass string) bool {
	return true
}
