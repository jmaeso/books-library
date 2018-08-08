package library

type User struct {
	Username string `db:"username"`
	Secret   []byte `db:"secret"`
}
