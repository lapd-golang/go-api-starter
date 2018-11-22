package models

import "go-admin-starter/database"

type User struct {
	Base

	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) CheckUser() (user User) {
	database.Eloquent.Where(u).First(&user)

	return
}
