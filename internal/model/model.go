package model

type UserPassword struct {
	ID       int    `db:"id"`
	Password string `db:"password"`
}
