package models

type AdminUser struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *AdminUser) CheckUser() (user AdminUser) {
	db.Where(u).First(&user)

	return
}
