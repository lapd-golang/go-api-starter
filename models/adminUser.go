package models

import "go-admin-starter/database"

type AdminUser struct {
	Base

	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *AdminUser) CheckUser() (user AdminUser) {
	database.Eloquent.Where(u).First(&user)

	return
}
