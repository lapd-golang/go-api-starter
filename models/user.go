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

func (u *User) CheckExistByUsername(username string) bool  {
	var user User
	db.Where("username = ?", username).First(&user)


	if user.ID > 0 {
		return true
	}

	return false
}

func (u *User) Insert() (id int, err error) {
	result := db.Create(&u)
	id = u.ID
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}
