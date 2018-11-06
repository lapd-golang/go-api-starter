package models

import "admin-server/database"

type User struct {
	Base

	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) GetByUsernameAndPassword(username string, password string) (user User) {
	database.Eloquent.Where(User{Username: username, Password: password}).First(&user)

	return
}
