package storage

type Connector struct {
	databaza string
	Users    []*Account
}

func NewConnector() *Connector {
	return &Connector{
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

func (st *Connector) CreateAccount(name, pass string) bool {
	user := &Account{
		username: name,
		password: pass,
	}
	st.Users = append(st.Users, user)
	return true
}

func (st *Connector) FindAccount(name, pass string) bool {
	for _, ak := range st.Users {
		if ak.username == name {
			return true
		}
	}
	return false
}
