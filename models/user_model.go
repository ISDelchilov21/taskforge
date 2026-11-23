package models

type User struct {
	ID           int
	Username     string
	Email        string
	Role         string
	PasswordHash string `json:"-"`
}
