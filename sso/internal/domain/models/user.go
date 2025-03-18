package models

type User struct {
	ID       int64
	UserName string
	Email    string
	PassHash []byte
	Role     string
}
