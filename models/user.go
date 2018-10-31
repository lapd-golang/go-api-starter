package models

type User struct {
	Base

	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) GetByUsernameAndPassword(username string, password string) (user User) {
	db.Where(User{Username: u.Username, Password: u.Password}).First(&user)

	return
}
