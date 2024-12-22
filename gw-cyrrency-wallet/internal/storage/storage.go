package storage

type Instance struct {
	databaza string
	Users    []*Account
	Ex       Exchange
}

func NewConnector() *Instance {
	return &Instance{
		databaza: "postgre",

		Users: []*Account{{
			username: "123zed",
			password: "qwerty",
		}, {
			username: "kartman",
			password: "south",
		}, {
			username: "bobr",
			password: "kzn",
		}},
	}
}

func (st *Instance) CreateAccount(name, pass string) bool {
	user := &Account{
		username: name,
		password: pass,
	}
	st.Users = append(st.Users, user)
	return true
}

func (st *Instance) FindAccount(name, pass string) bool {
	for _, ak := range st.Users {
		if ak.username == name && ak.password == pass {
			return true
		}
	}
	return false
}
