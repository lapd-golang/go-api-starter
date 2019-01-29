package models

type User struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) CheckUser() (user User) {
	db.Where(u).First(&user)

	return
}
