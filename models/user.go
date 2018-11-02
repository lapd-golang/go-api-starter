package models

type User struct {
	Base

	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) GetByUsernameAndPassword(username string, password string) (user User) {
	Eloquent.Where(User{Username: username, Password: password}).First(&user)

	return
}
