package models

type User struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) GetByUsername(username string) (user User)  {
	db.Where("username = ?", username).First(&user)
	return
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
